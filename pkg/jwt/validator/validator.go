package validator

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
	"slices"
	"time"
)

const defaultAllowedClockSkew = 0

var (
	// ErrKeyFuncMissing is returned when the validator is initialized without a
	// required key retrieval function.
	ErrKeyFuncMissing = errors.New("keyFunc is required (use WithKeyFunc)")

	// ErrSignatureAlgorithmMissing is returned when the validator is initialized
	// without required cryptographic signing algorithms.
	ErrSignatureAlgorithmMissing = errors.New("signature algorithm is required (use WithAlgorithm or WithAlgorithms)")

	// ErrIssuerMissing is returned when the validator is initialized without a
	// required issuer.
	ErrIssuerMissing = errors.New("issuer is required (use WithIssuer or WithIssuers)")

	// ErrAudienceMissing is returned when the validator is initialized without a
	// required audience.
	ErrAudienceMissing = errors.New("audience is required (use WithAudience or WithAudiences)")

	// ErrKeyFunc is returned when the key retrieval function fails to provide a key.
	ErrKeyFunc = errors.New("failed to get key from the key func")

	// ErrEmptyToken is returned when the token is empty.
	ErrEmptyToken = errors.New("token is empty")

	// ErrEmptyIssuerClaim is returned when the required 'iss' (issuer) claim is
	// missing from the token.
	ErrEmptyIssuerClaim = errors.New("token has no issuer claim (iss)")

	// ErrInvalidIssuerClaim is returned when the 'iss' claim does not match any of
	// the allowed or expected issuers.
	ErrInvalidIssuerClaim = errors.New("invalid issuer claim (iss)")

	// ErrInvalidAudienceClaim is returned when the 'aud' (audience) claim does not
	// match the expected recipient identifier.
	ErrInvalidAudienceClaim = errors.New("invalid audience claim (aud)")

	// ErrNotValidYet is returned when the current time is earlier than the token's
	// 'nbf' (not before) claim.
	ErrNotValidYet = errors.New("token not valid yet (nbf)")

	// ErrExpired is returned when the current time is past the token's 'exp'
	// (expiration time) claim.
	ErrExpired = errors.New("token is expired (exp)")

	// ErrIssuedInTheFuture is returned when the 'iat' (issued at) timestamp is set
	// ahead of the current server time.
	ErrIssuedInTheFuture = errors.New("token issued in the future (iat)")
)

// Validator defines the behavior for checking the integrity and authenticity of
// security tokens. Implementations of this interface are responsible for
// parsing, cryptographically verifying, and validating the claims of an incoming
// token string.
type Validator interface {
	// ValidateToken validates the passed token and returns the validated claims from
	// the token. It returns the extracted RegisteredClaims if the token is valid, or
	// an error if verification fails.
	ValidateToken(ctx context.Context, rawToken string) (claims.RegisteredClaims, error)
}

// validator is used to validate JWT.
type validator struct {
	keyFunc           func(context.Context) (any, error) // Required.
	allowedAlgorithms []crypto.SignatureAlgorithm        // Required.
	expectedAudiences []string                           // Required.
	expectedIssuers   []string                           // Required.
	allowedClockSkew  time.Duration                      // Optional.
}

// New returns new JWT validator instance.
func New(opts ...Option) (*validator, error) {
	v := &validator{
		allowedClockSkew: defaultAllowedClockSkew,
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, fmt.Errorf("invalid option: %w", err)
		}
	}

	// Validate required configuration
	if err := v.validateConfiguration(); err != nil {
		return nil, fmt.Errorf("invalid validator configuration: %w", err)
	}

	return v, nil
}

// validateConfiguration ensures all required fields are set.
func (v *validator) validateConfiguration() error {
	var errs []error

	if v.keyFunc == nil {
		errs = append(errs, ErrKeyFuncMissing)
	}

	if len(v.allowedAlgorithms) == 0 {
		errs = append(errs, ErrSignatureAlgorithmMissing)
	}

	if len(v.expectedIssuers) == 0 {
		errs = append(errs, ErrIssuerMissing)
	}

	if len(v.expectedAudiences) == 0 {
		errs = append(errs, ErrAudienceMissing)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// ValidateToken validates the passed token and returns the validated claims from
// the token. It returns the extracted RegisteredClaims if the token is valid, or
// an error if verification fails.
func (v *validator) ValidateToken(ctx context.Context, rawToken string) (claims.RegisteredClaims, error) {
	if rawToken == "" {
		return claims.RegisteredClaims{}, ErrEmptyToken
	}

	token, err := jwt.ParseToken(rawToken, v.allowedAlgorithms)
	if err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w", err)
	}

	unverifiedClaims, err := jwt.UnsafeParseTokenClaims(token)
	if err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w", err)
	}

	if err = v.validateIssuer(unverifiedClaims.Issuer); err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w", err)
	}

	key, err := v.keyFunc(ctx)
	if err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w: %w", ErrKeyFunc, err)
	}

	tokenClaims, err := jwt.ParseTokenClaims(token, key)
	if err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w", err)
	}

	if err = v.validateClaims(tokenClaims); err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("validation failed: %w", err)
	}

	return tokenClaims, nil
}

// validateIssuer checks if the token issuer matches one of the expected issuers.
func (v *validator) validateIssuer(issuer string) error {
	for _, expectedIssuer := range v.expectedIssuers {
		if issuer == expectedIssuer {
			return nil
		}
	}

	return fmt.Errorf("%w: token issuer %q does not match any expected issuer", ErrInvalidIssuerClaim, issuer)
}

// validateClaims checks claims in a token against expected values.
func (v *validator) validateClaims(claims claims.RegisteredClaims) error {
	leeway := v.allowedClockSkew
	expectedIssuers := v.expectedIssuers
	expectedAudiences := v.expectedAudiences

	if slices.Contains(expectedIssuers, claims.Issuer) == false {
		return ErrInvalidIssuerClaim
	}

	if extslices.Any(expectedAudiences, claims.Audience) == false {
		return ErrInvalidAudienceClaim
	}

	now := time.Now()
	if claims.NotBefore != nil && now.Add(leeway).Before(claims.NotBefore.Time()) {
		return ErrNotValidYet
	}

	if claims.Expiry != nil && now.Add(-leeway).After(claims.Expiry.Time()) {
		return ErrExpired
	}

	if claims.IssuedAt != nil && now.Add(leeway).Before(claims.IssuedAt.Time()) {
		return ErrIssuedInTheFuture
	}

	return nil
}
