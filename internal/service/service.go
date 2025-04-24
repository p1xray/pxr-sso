package service

import (
	"context"
	"pxr-sso/internal/dto"
)

// AuthService is service for working with user authentication and authorization.
type AuthService interface {
	// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
	Login(ctx context.Context, data *dto.LoginDTO) (*dto.TokensDTO, error)
}
