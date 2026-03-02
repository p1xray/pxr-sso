package entity

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
	"github.com/p1xray/pxr-sso/pkg/nullable"
)

// OAuthOption is how options for the OAuth are set up.
type OAuthOption func(*OAuth)

// WithBuilderURI is an option which sets up the URI builder for the OAuth.
func WithBuilderURI(uriBuilder *builder.URI) OAuthOption {
	return func(a *OAuth) {
		a.uriBuilder = uriBuilder
	}
}

// WithTokenGenerator is an option which sets up the token generator for the OAuth.
func WithTokenGenerator(tokenGenerator *generator.Token) OAuthOption {
	return func(a *OAuth) {
		a.tokenGenerator = tokenGenerator
	}
}

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

// WithAuthorization is an option which sets up the authorization for the OAuth.
func WithAuthorization(authorization dto.Authorization) OAuthOption {
	return func(a *OAuth) {
		a.authorization = nullable.Some(authorization)
	}
}
