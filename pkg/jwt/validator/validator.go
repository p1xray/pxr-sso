package validator

import (
	"context"
	"errors"
	"fmt"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	jwtparser "github.com/p1xray/pxr-sso/pkg/jwt/parser"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

var (
	ErrEmptyKeyFunc              = errors.New("keyFunc is required")
	ErrEmptyIssuer               = errors.New("issuer is required")
	ErrEmptyAudience             = errors.New("audience is required")
	ErrParsingToken              = errors.New("error parsing token")
	ErrParsingSignatureAlgorithm = errors.New("error parsing signature algorithm")
	ErrInvalidSigningMethod      = errors.New("signing method is invalid")
	ErrDeserializingTokenClaims  = errors.New("error deserializing token claims")
	ErrValidatingClaims          = errors.New("error validating token claims")
	ErrValidatingCustomClaims    = errors.New("error validating  token custom claims")
	ErrGettingKey                = errors.New("error getting key")
)

// Validator is used to validate JWT.
type Validator struct {
	keyFunc            func(context.Context) ([]byte, error)
	signatureAlgorithm jose.SignatureAlgorithm
	expectedClaims     jwt.Expected
	customClaims       func() jwtclaims.CustomClaims
	allowedClockSkew   time.Duration
}

// New returns new JWT validator instance.
func New(
	keyFunc func(context.Context) ([]byte, error),
	issuer string,
	audience []string,
	options ...Option,
) (*Validator, error) {
	if keyFunc == nil {
		return nil, ErrEmptyKeyFunc
	}

	if issuer == "" {
		return nil, ErrEmptyIssuer
	}

	if len(audience) == 0 {
		return nil, ErrEmptyAudience
	}

	validator := &Validator{
		keyFunc:            keyFunc,
		signatureAlgorithm: jose.HS256,
		expectedClaims: jwt.Expected{
			Issuer:      issuer,
			AnyAudience: audience,
		},
	}

	for _, opt := range options {
		opt(validator)
	}

	return validator, nil
}

// ValidateToken validates the passed token and returns the validated claims from the token.
func (v *Validator) ValidateToken(ctx context.Context, tokenString string) (jwtclaims.ValidatedClaims, error) {
	token, err := jwt.ParseSigned(tokenString, []jose.SignatureAlgorithm{v.signatureAlgorithm})
	if err != nil {
		return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrParsingToken, err)
	}

	signatureAlgorithm, err := jwtparser.ParseSignatureAlgorithm(token)
	if err != nil {
		return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrParsingSignatureAlgorithm, err)
	}

	if err = validateSigningMethod(v.signatureAlgorithm, signatureAlgorithm); err != nil {
		return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrInvalidSigningMethod, err)
	}

	registeredClaims, customClaims, err := v.parseClaims(ctx, token)
	if err != nil {
		return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrDeserializingTokenClaims, err)
	}

	if err = validateClaimsWithLeeway(registeredClaims.Claims, v.expectedClaims, v.allowedClockSkew); err != nil {
		return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrValidatingClaims, err)
	}

	if customClaims != nil {
		if err = customClaims.Validate(ctx); err != nil {
			return jwtclaims.ValidatedClaims{}, fmt.Errorf("%w: %w", ErrValidatingCustomClaims, err)
		}
	}

	validatedClaims := jwtclaims.ValidatedClaims{
		RegisteredClaims: registeredClaims,
		CustomClaims:     customClaims,
	}

	return validatedClaims, nil
}

func (v *Validator) parseClaims(
	ctx context.Context,
	token *jwt.JSONWebToken,
) (jwtclaims.AccessTokenClaims, jwtclaims.CustomClaims, error) {
	key, err := v.keyFunc(ctx)
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, nil, fmt.Errorf("%w: %w", ErrGettingKey, err)
	}

	registeredClaims, customClaims, err := jwtparser.ParseAccessToken(token, key, v.customClaims)
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, nil, fmt.Errorf("%w: %w", ErrParsingToken, err)
	}

	return registeredClaims, customClaims, nil
}

func validateSigningMethod(validAlgorithmName, tokenAlgorithmName jose.SignatureAlgorithm) error {
	if validAlgorithmName != tokenAlgorithmName {
		return fmt.Errorf("expected %q signing algorithm but token specified %q", validAlgorithmName, tokenAlgorithmName)
	}
	return nil
}

func validateClaimsWithLeeway(claims jwt.Claims, expected jwt.Expected, leeway time.Duration) error {
	expectedClaims := expected
	expectedClaims.Time = time.Now()

	if claims.Issuer != expectedClaims.Issuer {
		return jwt.ErrInvalidIssuer
	}

	isAudienceFound := false
	for _, aud := range claims.Audience {
		if expectedClaims.AnyAudience.Contains(aud) {
			isAudienceFound = true
			break
		}
	}

	if !isAudienceFound {
		return jwt.ErrInvalidAudience
	}

	if claims.NotBefore != nil && expectedClaims.Time.Add(leeway).Before(claims.NotBefore.Time()) {
		return jwt.ErrNotValidYet
	}

	if claims.Expiry != nil && expectedClaims.Time.Add(-leeway).After(claims.Expiry.Time()) {
		return jwt.ErrExpired
	}

	if claims.IssuedAt != nil && expectedClaims.Time.Add(leeway).Before(claims.IssuedAt.Time()) {
		return jwt.ErrIssuedInTheFuture
	}

	return nil
}
