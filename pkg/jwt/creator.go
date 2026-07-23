package jwt

import (
	"errors"
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
)

var (
	// ErrCreateSigner is returned when the cryptographic signing module fails to
	// initialize with the provided private key or algorithm.
	ErrCreateSigner = errors.New("failed to create signer")

	// ErrTokenSerialize is returned when the internal claims structure cannot be
	// marshaled into its final JSON or compact JWT string representation.
	ErrTokenSerialize = errors.New("failed to serialize token")
)

// CreateAccessToken creates a new signed JWT access token embedded with the
// provided claims.
func CreateAccessToken(claims claims.AccessTokenClaims, key []byte) (string, error) {
	token, err := createSignedTokenWithClaims(claims, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateRefreshToken creates a new signed JWT refresh token used for session
// renewal and obtaining new access/ID tokens.
func CreateRefreshToken(claims claims.RefreshTokenClaims, key []byte) (string, error) {
	token, err := createSignedTokenWithClaims(claims, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateIDToken creates a new signed JWT ID token embedded with the provided
// claims as defined in the OpenID Connect Core 1.0 specification.
func CreateIDToken(claims claims.IDTokenClaims, key []byte) (string, error) {
	token, err := createSignedTokenWithClaims(claims, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

// createSignedTokenWithClaims returns cryptographically signs a JWT, embedding
// the provided arbitrary claims.
func createSignedTokenWithClaims(claims any, key []byte) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCreateSigner, err)
	}

	tokenBuilder := jwt.Signed(sig)
	tokenBuilder = tokenBuilder.Claims(claims)

	tokenStr, err := tokenBuilder.Serialize()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenSerialize, err)
	}

	return tokenStr, nil
}
