package controller

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/login"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/logout"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/refresh"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/register"
)

type (
	// Login is a use-case for logging in a user.
	Login interface {
		// Execute executes the use-case for logging in a user. If successful, new tokens are returned.
		Execute(ctx context.Context, data login.Params) (entity.Tokens, error)
	}

	// Register is a use-case for registering a new user.
	Register interface {
		// Execute executes the use-case for registering a new user. If successful, new tokens are returned.
		Execute(ctx context.Context, data register.Params) (entity.Tokens, error)
	}

	// RefreshTokens is a use-case for refreshing user tokens.
	RefreshTokens interface {
		// Execute executes the use-case for refreshing user tokens. If successful, new tokens are returned.
		Execute(ctx context.Context, data refresh.Params) (entity.Tokens, error)
	}

	// Logout is a use-case for logging out a user.
	Logout interface {
		// Execute executes the use-case for logging out a user.
		Execute(ctx context.Context, data logout.Params) error
	}

	// UserProfile is a use-case for getting user profile data.
	UserProfile interface {
		// Execute executes the use-case for getting user profile data.
		Execute(ctx context.Context, id int64) (entity.User, error)
	}
)
