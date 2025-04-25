package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"pxr-sso/internal/domain"
	"pxr-sso/internal/dto"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/lib/token"
	"pxr-sso/internal/service"
	"pxr-sso/internal/storage"
	"time"
)

// Auth is service for working with user authentication and authorization.
type Auth struct {
	log             *slog.Logger
	storage         storage.SSOStorage
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// New creates a new auth service.
func New(
	log *slog.Logger,
) *Auth {
	return &Auth{log: log}
}

// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
func (a *Auth) Login(ctx context.Context, data *dto.LoginDTO) (*dto.TokensDTO, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
	)
	log.Info("attempting to login user")

	// Get user from storage.
	user, err := a.storage.UserByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return nil, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		a.log.Warn("invalid credentials", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
	}

	// Get user client by user link and client code from storage.
	client, err := a.storage.UserClient(ctx, user.ID, data.ClientCode)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create auth tokens.
	refreshToken := token.NewRefreshToken()
	accessToken, err := token.NewAccessToken(user, client, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := &domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.storage.CreateSession(ctx, sessionToCreate); err != nil {
		a.log.Error("failed to create session", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: returns tokens in response
	return &dto.TokensDTO{AccessToken: "", RefreshToken: ""}, nil
}
