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
)

// Logout terminates the user's session.
func (a *Auth) Logout(ctx context.Context, data dto.LogoutDTO) error {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
	)
	log.Info("attempting to user logout")

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, service.ErrClientNotFound)
		}

		log.Error("failed to get client", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := jwtparser.ParseRefreshToken(data.RefreshToken, []byte(client.SecretKey))
	if err != nil {
		log.Error("failed to parse refresh token", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Get session by refresh token from storage.
	session, err := a.sessionCRUD.SessionByRefreshToken(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("session not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, service.ErrSessionNotFound)
		}

		log.Error("failed to get session", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Remove current session from storage.
	if err = a.sessionCRUD.RemoveSession(ctx, session.ID); err != nil {
		log.Error("failed to remove session", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logout successfully")

	return nil
}
