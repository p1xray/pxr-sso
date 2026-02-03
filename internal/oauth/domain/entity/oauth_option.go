package entity

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/nullable"
)

// OAuthOption is how options for the OAuth are set up.
type OAuthOption func(*OAuth)

// WithNullableClient is an option which sets up the nullable client for the OAuth.
func WithNullableClient(client nullable.Nullable[dto.Client]) OAuthOption {
	return func(a *OAuth) {
		a.client = client
	}
}

// WithClient is an option which sets up the client for the OAuth.
func WithClient(client dto.Client) OAuthOption {
	return func(a *OAuth) {
		a.client = nullable.Some(client)
	}
}

// WithFlow is an option which sets up the flow for the OAuth.
func WithFlow(flow dto.Flow) OAuthOption {
	return func(a *OAuth) {
		a.flow = nullable.Some(flow)
	}
}

// WithUser is an option which sets up the user for the OAuth.
func WithUser(user dto.User) OAuthOption {
	return func(a *OAuth) {
		a.user = nullable.Some(user)
	}
}
