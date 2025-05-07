package service

import (
	"context"
	"errors"
	"pxr-sso/internal/dto"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserExists          = errors.New("user already exists")
	ErrSessionNotFound     = errors.New("session not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidSession      = errors.New("invalid session")
	ErrUserNotFound        = errors.New("user not found")
)

// AuthService is service for working with user authentication and authorization.
type AuthService interface {
	// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
	Login(ctx context.Context, data dto.LoginDTO) (dto.TokensDTO, error)

	// Register registers new user in the system and returns access and refresh tokens.
	Register(ctx context.Context, data dto.RegisterDTO) (dto.TokensDTO, error)

	// RefreshTokens refreshes the user's auth tokens.
	RefreshTokens(ctx context.Context, data dto.RefreshTokensDTO) (dto.TokensDTO, error)
}
