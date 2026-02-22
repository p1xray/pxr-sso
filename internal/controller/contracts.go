package controller

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/consent"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/register"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/token"
)

type (
	// Authorize is a use-case for OAuth authorize.
	Authorize interface {
		Execute(ctx context.Context, data authorize.Params) string
	}

	// Login is a use-case for logging in a user.
	Login interface {
		Execute(ctx context.Context, data login.Params) (string, error)
	}

	// Register is a use-case for registering a new user.
	Register interface {
		Execute(ctx context.Context, data register.Params) (string, error)
	}

	// Consent is a use-case for confirming consent.
	Consent interface {
		Execute(ctx context.Context, data consent.Params) (string, error)
	}

	// Token is a use-case for exchange token.
	Token interface {
		Execute(ctx context.Context, data token.Params) (dto.Token, error)
	}
)
