package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/logic/dto"
	"pxr-sso/internal/logic/service"
	"pxr-sso/internal/storage"
	jwtparser "pxr-sso/pkg/jwt/parser"
	"time"
)

// RefreshTokens refreshes the user's auth tokens.
func (a *Auth) RefreshTokens(ctx context.Context, data dto.RefreshTokensDTO) (dto.TokensDTO, error) {
	const op = "auth.RefreshTokens"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
	)
	log.Info("attempting to refresh tokens")

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrClientNotFound)
		}

		log.Error("failed to get client", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := jwtparser.ParseRefreshToken(data.RefreshToken, []byte(client.SecretKey))
	if err != nil {
		log.Error("failed to parse refresh token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get session by refresh token ID from storage.
	session, err := a.sessionCRUD.SessionByRefreshToken(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("session not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrSessionNotFound)
		}

		log.Error("failed to get session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check session expiration time.
	now := time.Now()
	if session.ExpiresAt.Before(now) {
		log.Warn("refresh token expired")

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrRefreshTokenExpired)
	}

	// Check session user agent and fingerprint.
	if session.UserAgent != data.UserAgent && session.Fingerprint != data.Fingerprint {
		log.Warn("invalid session")

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidSession)
	}

	// Get user from storage.
	user, err := a.userCRUD.UserWithPermission(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Remove current session from storage.
	if err = a.sessionCRUD.RemoveSession(ctx, session.ID); err != nil {
		log.Error("failed to remove session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access and refresh tokens.
	tokensData, err := a.createAccessAndRefreshTokens(log, user, client, data.Issuer)
	if err != nil {
		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := dto.CreateSessionDTO{
		UserID:       user.ID,
		RefreshToken: tokensData.RefreshTokenID,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.sessionCRUD.CreateSession(ctx, sessionToCreate); err != nil {
		log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("tokens refreshed successfully")

	return dto.TokensDTO{AccessToken: tokensData.AccessToken, RefreshToken: tokensData.RefreshToken}, nil
}
