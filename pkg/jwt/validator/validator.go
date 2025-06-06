package validator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/go-jose/go-jose.v2/jwt"
)

// Signature algorithms
const (
	EdDSA = SignatureAlgorithm("EdDSA")
	HS256 = SignatureAlgorithm("HS256") // HMAC using SHA-256
	HS384 = SignatureAlgorithm("HS384") // HMAC using SHA-384
	HS512 = SignatureAlgorithm("HS512") // HMAC using SHA-512
	RS256 = SignatureAlgorithm("RS256") // RSASSA-PKCS-v1.5 using SHA-256
	RS384 = SignatureAlgorithm("RS384") // RSASSA-PKCS-v1.5 using SHA-384
	RS512 = SignatureAlgorithm("RS512") // RSASSA-PKCS-v1.5 using SHA-512
	ES256 = SignatureAlgorithm("ES256") // ECDSA using P-256 and SHA-256
	ES384 = SignatureAlgorithm("ES384") // ECDSA using P-384 and SHA-384
	ES512 = SignatureAlgorithm("ES512") // ECDSA using P-521 and SHA-512
	PS256 = SignatureAlgorithm("PS256") // RSASSA-PSS using SHA256 and MGF1-SHA256
	PS384 = SignatureAlgorithm("PS384") // RSASSA-PSS using SHA384 and MGF1-SHA384
	PS512 = SignatureAlgorithm("PS512") // RSASSA-PSS using SHA512 and MGF1-SHA512
)

type Validator struct {
	keyFunc            func(context.Context) ([]byte, error)
	signatureAlgorithm SignatureAlgorithm
	expectedClaims     jwt.Expected
	customClaims       func() CustomClaims
	allowedClockSkew   time.Duration
}

// SignatureAlgorithm represents a signature (or MAC) algorithm.
type SignatureAlgorithm string

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
		signatureAlgorithm: HS256,
		expectedClaims: jwt.Expected{
			Issuer:   issuer,
			Audience: audience,
		},
	}

	for _, opt := range options {
		opt(validator)
	}

	return validator, nil
}

func (v *Validator) ValidateToken(ctx context.Context, tokenString string) (ValidatedClaims, error) {
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return ValidatedClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	if err = validateSigningMethod(v.signatureAlgorithm, SignatureAlgorithm(token.Headers[0].Algorithm)); err != nil {
		return ValidatedClaims{}, fmt.Errorf("signing method is invalid: %w", err)
	}

	registeredClaims, customClaims, err := v.deserializeClaims(ctx, token)
	if err != nil {
		return ValidatedClaims{}, fmt.Errorf("error deserializing token claims: %w", err)
	}

	if err = validateClaimsWithLeeway(registeredClaims, v.expectedClaims, v.allowedClockSkew); err != nil {
		return ValidatedClaims{}, fmt.Errorf("error validating claims: %w", err)
	}

	if customClaims != nil {
		if err = customClaims.Validate(ctx); err != nil {
			return ValidatedClaims{}, fmt.Errorf("error validating custom claims: %w", err)
		}
	}

	validatedClaims := ValidatedClaims{
		RegisteredClaims: RegisteredClaims{
			Issuer:    registeredClaims.Issuer,
			Subject:   registeredClaims.Subject,
			Audience:  registeredClaims.Audience,
			Expiry:    numericDateToUnixTime(registeredClaims.Expiry),
			NotBefore: numericDateToUnixTime(registeredClaims.NotBefore),
			IssuedAt:  numericDateToUnixTime(registeredClaims.IssuedAt),
			ID:        registeredClaims.ID,
		},
		CustomClaims: customClaims,
	}

	return validatedClaims, nil
}

func (v *Validator) deserializeClaims(ctx context.Context, token *jwt.JSONWebToken) (jwt.Claims, CustomClaims, error) {
	key, err := v.keyFunc(ctx)
	if err != nil {
		return jwt.Claims{}, nil, fmt.Errorf("error getting key: %w", err)
	}

	claims := []interface{}{&jwt.Claims{}}
	if v.customClaimsExist() {
		claims = append(claims, v.customClaims())
	}

	if err = token.Claims(key, claims...); err != nil {
		return jwt.Claims{}, nil, fmt.Errorf("error getting token claims: %w", err)
	}

	registeredClaims := *claims[0].(*jwt.Claims)

	var customClaims CustomClaims
	if len(claims) > 1 {
		customClaims = claims[1].(CustomClaims)
	}

	return registeredClaims, customClaims, nil
}

func (v *Validator) customClaimsExist() bool {
	return v.customClaims != nil && v.customClaims() != nil
}

func validateSigningMethod(validAlgorithmName, tokenAlgorithmName SignatureAlgorithm) error {
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

func numericDateToUnixTime(date *jwt.NumericDate) int64 {
	if date != nil {
		return date.Time().Unix()
	}
	return 0
}
