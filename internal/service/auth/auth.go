package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"pxr-sso/internal/dto"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/service"
	"pxr-sso/internal/storage"
)

// Auth is service for working with user authentication and authorization.
type Auth struct {
	log     *slog.Logger
	storage storage.SSOStorage
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

	user, err := a.storage.UserByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return nil, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		a.log.Warn("invalid credentials", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
	}

	// TODO: get user client by user link and client_code

	// TODO: create auth tokens

	// TODO: returns tokens in response
	return &dto.TokensDTO{AccessToken: "", RefreshToken: ""}, nil
}
