package service

import (
	"context"
	"errors"
	"pxr-sso/internal/dto"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService is service for working with user authentication and authorization.
type AuthService interface {
	// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
	Login(ctx context.Context, data dto.LoginDTO) (dto.TokensDTO, error)
}
