package jwtclaims

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
	TokenType string `json:"token_type,omitempty"`
	Scope     string `json:"scope,omitempty"`
}

// AccessTokenClaims are all access token claims of the current SSO project.
type AccessTokenClaims struct {
	jwt.Claims
	RegisteredCustomClaims
}

// RefreshTokenClaims are refresh token claims of the current SSO project.
type RefreshTokenClaims struct {
	ID        string           `json:"jti,omitempty"`
	TokenType string           `json:"token_type,omitempty"`
	Expiry    *jwt.NumericDate `json:"exp,omitempty"`
}

// CustomClaims defines any custom data / claims wanted.
type CustomClaims interface {
	Validate(ctx context.Context) error
}

func NumericDateToInt64(v *jwt.NumericDate) int64 {
	if v == nil {
		return 0
	}

	return int64(*v)
}
