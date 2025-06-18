package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	"github.com/p1xray/pxr-sso/internal/logic/service"
	"github.com/p1xray/pxr-sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
func (a *Auth) Login(ctx context.Context, data dto.LoginDTO) (dto.TokensDTO, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
	)
	log.Info("attempting to login user")

	// Get user data from storage.
	user, err := a.userCRUD.UserWithPermissionsByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check password hash.
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		log.Warn("invalid credentials", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
	}

	// Get user client by user link and client code from storage.
	client, err := a.clientCRUD.UserClient(ctx, user.ID, data.ClientCode)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrClientNotFound)
		}

		log.Error("failed to get client", sl.Err(err))

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

	log.Info("user logged in successfully")

	return dto.TokensDTO{AccessToken: tokensData.AccessToken, RefreshToken: tokensData.RefreshToken}, nil
}
