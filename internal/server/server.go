package server

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthService is service for working with user authentication and authorization.
type AuthService interface {
	// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
	Login(ctx context.Context, data dto.LoginDTO) (dto.TokensDTO, error)

	// Register registers new user in the system and returns access and refresh tokens.
	Register(ctx context.Context, data dto.RegisterDTO) (dto.TokensDTO, error)

	// RefreshTokens refreshes the user's auth tokens.
	RefreshTokens(ctx context.Context, data dto.RefreshTokensDTO) (dto.TokensDTO, error)

	// Logout terminates the user's session.
	Logout(ctx context.Context, data dto.LogoutDTO) error
}

// ProfileService is service for working with user profile.
type ProfileService interface {
	// UserProfile returns user profile data.
	UserProfile(ctx context.Context, userID int64) (dto.UserProfileDTO, error)
}

func InvalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func InternalError(msg string) error {
	return status.Error(codes.Internal, msg)
}

func NotFoundError(msg string) error {
	return status.Error(codes.NotFound, msg)
}
