package entity

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/nullable"
)

// OAuthOption is how options for the OAuth are set up.
type OAuthOption func(*OAuth)

// WithClient is an option which sets up the client for the OAuth.
func WithClient(client nullable.Nullable[dto.Client]) OAuthOption {
	return func(a *OAuth) {
		a.client = client
	}
}

// WithFlow is an option which sets up the flow for the OAuth.
func WithFlow(flow nullable.Nullable[dto.Flow]) OAuthOption {
	return func(a *OAuth) {
		a.flow = flow
	}
}

// WithUser is an option which sets up the user for the OAuth.
func WithUser(user nullable.Nullable[dto.User]) OAuthOption {
	return func(a *OAuth) {
		a.user = user
	}
}
