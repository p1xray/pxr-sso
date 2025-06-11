package jwtmiddleware

import (
	"context"
	"github.com/go-jose/go-jose/v4/jwt"
)

// ValidatedClaims is the struct that will be inserted into the context for the user.
type ValidatedClaims struct {
	RegisteredClaims AccessTokenClaims
	CustomClaims     CustomClaims
}

// RegisteredCustomClaims are custom claims of the current SSO project.
type RegisteredCustomClaims struct {
	Scope string `json:"scope,omitempty"`
}

// AccessTokenClaims are all access token claims of the current SSO project.
type AccessTokenClaims struct {
	jwt.Claims
	RegisteredCustomClaims
}

// RefreshTokenClaims are refresh token claims of the current SSO project.
type RefreshTokenClaims struct {
	ID     string           `json:"jti,omitempty"`
	Expiry *jwt.NumericDate `json:"exp,omitempty"`
}

// CustomClaims defines any custom data / claims wanted.
type CustomClaims interface {
	Validate(ctx context.Context) error
}
