package validator

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/pkg/jwt"
	jwtparser "github.com/p1xray/pxr-sso/pkg/jwt/parser"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// Validator is used to validate JWT.
type Validator struct {
	keyFunc            func(context.Context) ([]byte, error)
	signatureAlgorithm jose.SignatureAlgorithm
	expectedClaims     jwt.Expected
	customClaims       func() jwtmiddleware.CustomClaims
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
		return nil, errors.New("keyFunc is required")
	}

	if issuer == "" {
		return nil, errors.New("issuer is required")
	}

	if len(audience) == 0 {
		return nil, errors.New("audience is required")
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
func (v *Validator) ValidateToken(ctx context.Context, tokenString string) (jwtmiddleware.ValidatedClaims, error) {
	token, err := jwt.ParseSigned(tokenString, []jose.SignatureAlgorithm{v.signatureAlgorithm})
	if err != nil {
		return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	signatureAlgorithm, err := jwtparser.ParseSignatureAlgorithm(token)
	if err != nil {
		return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("error parsing signature algorithm: %w", err)
	}

	if err = validateSigningMethod(v.signatureAlgorithm, signatureAlgorithm); err != nil {
		return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("signing method is invalid: %w", err)
	}

	registeredClaims, customClaims, err := v.parseClaims(ctx, token)
	if err != nil {
		return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("error deserializing token claims: %w", err)
	}

	if err = validateClaimsWithLeeway(registeredClaims.Claims, v.expectedClaims, v.allowedClockSkew); err != nil {
		return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("error validating claims: %w", err)
	}

	if customClaims != nil {
		if err = customClaims.Validate(ctx); err != nil {
			return jwtmiddleware.ValidatedClaims{}, fmt.Errorf("error validating custom claims: %w", err)
		}
	}

	validatedClaims := jwtmiddleware.ValidatedClaims{
		RegisteredClaims: registeredClaims,
		CustomClaims:     customClaims,
	}

	return validatedClaims, nil
}

func (v *Validator) parseClaims(
	ctx context.Context,
	token *jwt.JSONWebToken,
) (jwtmiddleware.AccessTokenClaims, jwtmiddleware.CustomClaims, error) {
	key, err := v.keyFunc(ctx)
	if err != nil {
		return jwtmiddleware.AccessTokenClaims{}, nil, fmt.Errorf("error getting key: %w", err)
	}

	registeredClaims, customClaims, err := jwtparser.ParseAccessToken(token, key, v.customClaims)
	if err != nil {
		return jwtmiddleware.AccessTokenClaims{}, nil, fmt.Errorf("error parsing token: %w", err)
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
		if claims.Audience.Contains(aud) {
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
