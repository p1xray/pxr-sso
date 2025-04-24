package auth

import (
	"context"
	"log/slog"
	"pxr-sso/internal/dto"
)

// Auth is service for working with user authentication and authorization.
type Auth struct {
	log *slog.Logger
}

// New creates a new auth service.
func New(
	log *slog.Logger,
) *Auth {
	return &Auth{log: log}
}

// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
func (a *Auth) Login(ctx context.Context, data *dto.LoginDTO) (*dto.TokensDTO, error) {

	// TODO: get user from storage

	// TODO: check password hash

	// TODO: get user client by user link and client_code

	// TODO: create auth tokens

	// TODO: returns tokens in response
	return &dto.TokensDTO{AccessToken: "", RefreshToken: ""}, nil
}
