package validator

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
	"net/url"
	"time"
)

var (
	ErrEmptyKeyFunc           = errors.New("keyFunc cannot be nil")
	ErrUseBothAlgorithmOption = errors.New("cannot use WithAlgorithm with WithAlgorithms")
	ErrUnsupportedAlgorithm   = errors.New("unsupported signature algorithm")
	ErrEmptyAlgorithms        = errors.New("algorithms list cannot be empty")
	ErrEmptyIssuer            = errors.New("issuer cannot be empty")
	ErrEmptyIssuers           = errors.New("issuers list cannot be empty")
	ErrInvalidIssuerURL       = errors.New("invalid issuer URL")
	ErrEmptyAudience          = errors.New("audience cannot be empty")
	ErrEmptyAudiences         = errors.New("audiences list cannot be empty")
	ErrNegativeClockSkew      = errors.New("clock skew cannot be negative")
)

var allowedSignatureAlgorithms = map[crypto.SignatureAlgorithm]bool{
	crypto.EdDSA: true,
	crypto.HS256: true,
	crypto.HS384: true,
	crypto.HS512: true,
	crypto.RS256: true,
	crypto.RS384: true,
	crypto.RS512: true,
	crypto.ES256: true,
	crypto.ES384: true,
	crypto.ES512: true,
	crypto.PS256: true,
	crypto.PS384: true,
	crypto.PS512: true,
}

// Option is how options for the Validator are set up.
// Options return errors to enable validation during construction.
type Option func(*validator) error

// WithKeyFunc sets the function that provides the key for token verification.
// This is a required option.
//
// The keyFunc is called during token validation to retrieve the key(s) used
// to verify the token signature. For JWKS-based validation, use jwks.Provider.KeyFunc.
func WithKeyFunc(keyFunc func(context.Context) (any, error)) Option {
	return func(v *validator) error {
		if keyFunc == nil {
			return ErrEmptyKeyFunc
		}

		v.keyFunc = keyFunc
		return nil
	}
}

// WithAlgorithm sets the signature algorithm that tokens must use.
// This is a required option.
//
// Supported algorithms: EdDSA, HS256, HS384, HS512, RS256, RS384, RS512, ES256,
// ES384, ES512, PS256, PS384,PS512.
func WithAlgorithm(algorithm crypto.SignatureAlgorithm) Option {
	return func(v *validator) error {
		if len(v.allowedAlgorithms) > 0 {
			return ErrUseBothAlgorithmOption
		}

		if _, ok := allowedSignatureAlgorithms[algorithm]; !ok {
			return fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
		}

		v.allowedAlgorithms = []crypto.SignatureAlgorithm{algorithm}
		return nil
	}
}

// WithAlgorithms sets multiple signature algorithms that tokens may use.
// This is useful for mixed-algorithm MCD (Multiple Custom Domains) scenarios
// where different issuers use different algorithms (e.g., RS256 + HS256).
//
// Tokens will be validated against the allowed list before signature verification.
// Cannot be used with WithAlgorithm — they are mutually exclusive.
//
// Supported algorithms: EdDSA, HS256, HS384, HS512, RS256, RS384, RS512, ES256,
// ES384, ES512, PS256, PS384,PS512.
func WithAlgorithms(algorithms []crypto.SignatureAlgorithm) Option {
	return func(v *validator) error {
		if len(v.allowedAlgorithms) > 0 {
			return ErrUseBothAlgorithmOption
		}

		if len(algorithms) == 0 {
			return ErrEmptyAlgorithms
		}

		for i, alg := range algorithms {
			if _, ok := allowedSignatureAlgorithms[alg]; !ok {
				return fmt.Errorf("%w at index %d: %s", ErrUnsupportedAlgorithm, i, alg)
			}
		}

		v.allowedAlgorithms = algorithms
		return nil
	}
}

// WithIssuer sets a single expected issuer claim (iss) for token validation.
// This is a required option (use either WithIssuer, WithIssuers, or WithIssuersResolver, not multiple).
//
// The issuer URL should match the iss claim in the JWT. Tokens with a
// different issuer will be rejected.
//
// Cannot be used with WithIssuersResolver.
func WithIssuer(issuerURL string) Option {
	return func(v *validator) error {
		if issuerURL == "" {
			return ErrEmptyIssuer
		}

		if _, err := url.Parse(issuerURL); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidIssuerURL, err)
		}

		v.expectedIssuers = []string{issuerURL}
		return nil
	}
}

// WithIssuers sets multiple expected issuer claims (iss) for token validation.
// This is a required option (use either WithIssuer, WithIssuers, or WithIssuersResolver, not multiple).
//
// The token must contain one of the specified issuers. Tokens without
// any matching issuer will be rejected.
//
// Cannot be used with WithIssuersResolver.
func WithIssuers(issuers []string) Option {
	return func(v *validator) error {
		if len(issuers) == 0 {
			return ErrEmptyIssuers
		}

		for i, iss := range issuers {
			if iss == "" {
				return fmt.Errorf("%w at index %d", ErrEmptyIssuer, i)
			}

			if _, err := url.Parse(iss); err != nil {
				return fmt.Errorf("%w at index %d: %w", ErrInvalidIssuerURL, i, err)
			}
		}

		v.expectedIssuers = issuers
		return nil
	}
}

// WithAudience sets a single expected audience claim (aud) for token validation.
// This is a required option (use either WithAudience or WithAudiences, not both).
//
// The audience should match one of the aud claims in the JWT. Tokens without
// a matching audience will be rejected.
func WithAudience(audience string) Option {
	return func(v *validator) error {
		if audience == "" {
			return ErrEmptyAudience
		}

		v.expectedAudiences = []string{audience}
		return nil
	}
}

// WithAudiences sets multiple expected audience claims (aud) for token validation.
// This is a required option (use either WithAudience or WithAudiences, not both).
//
// The token must contain at least one of the specified audiences. Tokens without
// any matching audience will be rejected.
func WithAudiences(audiences []string) Option {
	return func(v *validator) error {
		if len(audiences) == 0 {
			return ErrEmptyAudiences
		}

		for i, aud := range audiences {
			if aud == "" {
				return fmt.Errorf("%w at index %d", ErrEmptyAudience, i)
			}
		}

		v.expectedAudiences = audiences
		return nil
	}
}

// WithAllowedClockSkew sets the allowed clock skew for time-based claims.
//
// This allows for some tolerance when validating exp, nbf, and iat claims
// to account for clock differences between systems. If not set, the default
// is 0 (no clock skew allowed).
func WithAllowedClockSkew(skew time.Duration) Option {
	return func(v *validator) error {
		if skew < 0 {
			return ErrNegativeClockSkew
		}

		v.allowedClockSkew = skew
		return nil
	}
}
