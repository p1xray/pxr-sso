package jwtcreator

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	jwtmiddleware "github.com/p1xray/pxr-sso/pkg/jwt"
	"strings"
	"time"
)

// AccessTokenCreateData is data to create new access token.
type AccessTokenCreateData struct {
	Subject  string
	Audience string
	Scopes   []string
	Issuer   string
	TTL      time.Duration
	Key      []byte
}

// NewAccessToken returns new JWT with claims.
func NewAccessToken(data AccessTokenCreateData) (string, error) {
	now := time.Now()
	claims := jwtmiddleware.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:        uuid.New().String(),
			Subject:   data.Subject,
			Issuer:    data.Issuer,
			Audience:  []string{data.Audience},
			Expiry:    jwt.NewNumericDate(now.Add(data.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		RegisteredCustomClaims: jwtmiddleware.RegisteredCustomClaims{
			Scope: strings.Join(data.Scopes, " "),
		},
	}

	tokenStr, err := createSignedTokenWithClaims(claims, data.Key)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// NewRefreshToken returns new refresh token.
func NewRefreshToken(key []byte, ttl time.Duration) (refreshToken string, refreshTokenID string, err error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtmiddleware.RefreshTokenClaims{
		ID:     id,
		Expiry: jwt.NewNumericDate(now.Add(ttl)),
	}

	tokenStr, err := createSignedTokenWithClaims(claims, key)
	if err != nil {
		return "", "", err
	}

	return tokenStr, id, nil
}

func createSignedTokenWithClaims(claims interface{}, key []byte) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	tokenStr, err := jwt.Signed(sig).Claims(claims).Serialize()
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
