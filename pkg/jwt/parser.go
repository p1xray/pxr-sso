package jwt

import (
	"errors"
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
)

var (
	// ErrParseToken indicates that the raw token string could not be parsed, either
	// due to an invalid format, malformed structure, or an unapproved signature algorithm.
	ErrParseToken = errors.New("failed to parse token")

	// ErrParseClaims indicates that the claims could not be extracted or unmarshaled
	// from the token payload, typically due to malformed JSON or invalid data types.
	ErrParseClaims = errors.New("failed to parse claims")
)

// ParseToken parses the raw JWT string and validates its structure against a whitelist of
// allowed cryptographic signature algorithms. It returns a decoded JSONWebToken object
// ready for subsequent signature verification and claims extraction, mitigating JWT algorithm
// confusion vulnerabilities.
func ParseToken(rawToken string, allowedSignatureAlgorithms []crypto.SignatureAlgorithm) (*jwt.JSONWebToken, error) {
	algorithms := make([]jose.SignatureAlgorithm, len(allowedSignatureAlgorithms))
	for i, alg := range allowedSignatureAlgorithms {
		algorithms[i] = jose.SignatureAlgorithm(alg)
	}

	token, err := jwt.ParseSigned(rawToken, algorithms)
	if err != nil {
		return &jwt.JSONWebToken{}, fmt.Errorf("%w: %w", ErrParseToken, err)
	}

	return token, nil
}

// UnsafeParseTokenClaims extracts the standard registered claims from the token
// without verifying its cryptographic signature.
//
// WARNING: This function is unsafe and should only be used when the token signature
// has already been validated, or when reading non-sensitive header data (e.g., token hints).
func UnsafeParseTokenClaims(token *jwt.JSONWebToken) (claims.RegisteredClaims, error) {
	registeredClaims := claims.RegisteredClaims{}
	if err := token.UnsafeClaimsWithoutVerification(&registeredClaims); err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("%w: %w", ErrParseClaims, err)
	}

	return registeredClaims, nil
}

// ParseTokenClaims parses token using a secret key into a set of registered claims.
func ParseTokenClaims(token *jwt.JSONWebToken, key any) (claims.RegisteredClaims, error) {
	registeredClaims := claims.RegisteredClaims{}
	if err := token.Claims(key, &registeredClaims); err != nil {
		return claims.RegisteredClaims{}, fmt.Errorf("%w: %w", ErrParseClaims, err)
	}

	return registeredClaims, nil
}

// ParseAccessTokenClaims parses access token using a secret key into a set of claims.
func ParseAccessTokenClaims(token *jwt.JSONWebToken, key any) (claims.AccessTokenClaims, error) {
	tokenClaims := claims.AccessTokenClaims{}
	if err := token.Claims(key, &tokenClaims); err != nil {
		return claims.AccessTokenClaims{}, fmt.Errorf("%w: %w", ErrParseClaims, err)
	}

	return tokenClaims, nil
}

// ParseRefreshTokenClaims parses refresh token using a secret key into a set of claims.
func ParseRefreshTokenClaims(token *jwt.JSONWebToken, key any) (claims.RefreshTokenClaims, error) {
	tokenClaims := claims.RefreshTokenClaims{}
	if err := token.Claims(key, &tokenClaims); err != nil {
		return claims.RefreshTokenClaims{}, fmt.Errorf("%w: %w", ErrParseClaims, err)
	}

	return tokenClaims, nil
}

// ParseIDTokenClaims parses id token using a secret key into a set of claims.
func ParseIDTokenClaims(token *jwt.JSONWebToken, key any) (claims.IDTokenClaims, error) {
	tokenClaims := claims.IDTokenClaims{}
	if err := token.Claims(key, &tokenClaims); err != nil {
		return claims.IDTokenClaims{}, fmt.Errorf("%w: %w", ErrParseClaims, err)
	}

	return tokenClaims, nil
}
