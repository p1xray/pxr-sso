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
	Login interface {
		Execute(ctx context.Context, data login.Params) (entity.Tokens, error)
	}

	Register interface {
		Execute(ctx context.Context, data register.Params) (entity.Tokens, error)
	}

	RefreshTokens interface {
		Execute(ctx context.Context, data refresh.Params) (entity.Tokens, error)
	}

	Logout interface {
		Execute(ctx context.Context, data logout.Params) error
	}

	UserProfile interface {
		Execute(ctx context.Context, id int64) (entity.User, error)
	}
)
