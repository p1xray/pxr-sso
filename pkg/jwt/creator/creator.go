package jwtcreator

import (
	"errors"
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"strings"
	"time"
)

var (
	ErrCreateSigner   = errors.New("error creating signer")
	ErrTokenSerialize = errors.New("error serializing token")
)

// AccessTokenCreateData is data to create new access token.
type AccessTokenCreateData struct {
	Subject      string
	Audiences    []string
	Scopes       []string
	Issuer       string
	CustomClaims map[string]interface{}
	TTL          time.Duration
	Key          []byte
}

// NewAccessToken returns new JWT with claims.
func NewAccessToken(data AccessTokenCreateData) (string, error) {
	now := time.Now()
	registeredClaims := jwtclaims.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:        uuid.New().String(),
			Subject:   data.Subject,
			Issuer:    data.Issuer,
			Audience:  data.Audiences,
			Expiry:    jwt.NewNumericDate(now.Add(data.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
			Scope: strings.Join(data.Scopes, " "),
		},
	}

	tokenStr, err := createSignedTokenWithClaims(data.Key, registeredClaims, data.CustomClaims)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// NewRefreshToken returns new refresh token.
func NewRefreshToken(key []byte, ttl time.Duration) (refreshToken string, refreshTokenID string, err error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtclaims.RefreshTokenClaims{
		ID:     id,
		Expiry: jwt.NewNumericDate(now.Add(ttl)),
	}

	tokenStr, err := createSignedTokenWithClaims(key, claims, nil)
	if err != nil {
		return "", "", err
	}

	return tokenStr, id, nil
}

func createSignedTokenWithClaims(key []byte, registeredClaims interface{}, customClaims interface{}) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCreateSigner, err)
	}

	tokenBuilder := jwt.Signed(sig)
	tokenBuilder = tokenBuilder.Claims(registeredClaims)

	if customClaims != nil {
		tokenBuilder = tokenBuilder.Claims(customClaims)
	}

	tokenStr, err := tokenBuilder.Serialize()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenSerialize, err)
	}

	return tokenStr, nil
}
