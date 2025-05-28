package validator

import (
	"context"
)

// ValidatedClaims is the struct that will be inserted into the context for the user.
type ValidatedClaims struct {
	RegisteredClaims RegisteredClaims
	CustomClaims     CustomClaims
}

// RegisteredClaims represents public claim values (as specified in RFC 7519).
type RegisteredClaims struct {
	Issuer    string   `json:"iss,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	Expiry    int64    `json:"exp,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	ID        string   `json:"jti,omitempty"`
}

// CustomClaims defines any custom data / claims wanted.
type CustomClaims interface {
	Validate(ctx context.Context) error
}
