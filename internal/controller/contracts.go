package controller

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/token"
	oldLogin "github.com/p1xray/pxr-sso/internal/usecase/auth/login"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/logout"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/refresh"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/register"
	"github.com/p1xray/pxr-sso/internal/usecase/profile/edit"
)

type (
	// OldLogin is a use-case for logging in a user.
	OldLogin interface {
		// Execute executes the use-case for logging in a user. If successful, new tokens are returned.
		Execute(ctx context.Context, data oldLogin.Params) (entity.Tokens, error)
	}

	// OldRegister is a use-case for registering a new user.
	OldRegister interface {
		// Execute executes the use-case for registering a new user. If successful, new tokens are returned.
		Execute(ctx context.Context, data register.Params) (entity.Tokens, error)
	}

	// OldRefreshTokens is a use-case for refreshing user tokens.
	OldRefreshTokens interface {
		// Execute executes the use-case for refreshing user tokens. If successful, new tokens are returned.
		Execute(ctx context.Context, data refresh.Params) (entity.Tokens, error)
	}

	// OldLogout is a use-case for logging out a user.
	OldLogout interface {
		// Execute executes the use-case for logging out a user.
		Execute(ctx context.Context, data logout.Params) error
	}

	// UserProfile is a use-case for getting user profile data.
	UserProfile interface {
		// Execute executes the use-case for getting user profile data.
		Execute(ctx context.Context, id int64) (entity.User, error)
	}

	// EditProfile is a use-case for editing user profile data.
	EditProfile interface {
		Execute(ctx context.Context, data edit.Params) error
	}

	// Authorize is a use-case for OAuth authorize.
	Authorize interface {
		Execute(ctx context.Context, data authorize.Params) string
	}

	// Login is a use-case for logging in a user.
	Login interface {
		Execute(ctx context.Context, data login.Params) (string, *domain.DisplayableError)
	}

	// Token is a use-case for exchange token.
	Token interface {
		Execute(ctx context.Context, data token.Params) (dto.Token, error)
	}
)
