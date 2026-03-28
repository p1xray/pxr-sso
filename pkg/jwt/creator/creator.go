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

// IDTokenCreateData is data to create new ID token.
type IDTokenCreateData struct {
	Subject   string
	ClientID  string
	Issuer    string
	AuthTime  time.Time
	Username  string
	Name      string
	Gender    string
	Birthdate time.Time
	Picture   string
	TTL       time.Duration
	Key       []byte
}

// NewAccessToken returns new JWT with claims.
func NewAccessToken(data AccessTokenCreateData) (jwtclaims.AccessTokenClaims, string, error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtclaims.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:        id,
			Subject:   data.Subject,
			Issuer:    data.Issuer,
			Audience:  data.Audiences,
			Expiry:    jwt.NewNumericDate(now.Add(data.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
			TokenType: "Bearer",
			Scope:     strings.Join(data.Scopes, " "),
		},
	}

	token, err := createSignedTokenWithClaims(data.Key, claims, data.CustomClaims)
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, "", err
	}

	return claims, token, nil
}

// NewRefreshToken returns new refresh token.
func NewRefreshToken(key []byte, ttl time.Duration) (jwtclaims.RefreshTokenClaims, string, error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtclaims.RefreshTokenClaims{
		ID:        id,
		TokenType: "refresh",
		Expiry:    jwt.NewNumericDate(now.Add(ttl)),
	}

	token, err := createSignedTokenWithClaims(key, claims, nil)
	if err != nil {
		return jwtclaims.RefreshTokenClaims{}, "", err
	}

	return claims, token, nil
}

// NewIDToken returns new ID token.
func NewIDToken(data IDTokenCreateData) (jwtclaims.IDTokenClaims, string, error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtclaims.IDTokenClaims{
		Claims: jwt.Claims{
			ID:        id,
			Subject:   data.Subject,
			Issuer:    data.Issuer,
			Audience:  []string{data.ClientID},
			Expiry:    jwt.NewNumericDate(now.Add(data.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		TokenType: "id_token",
		AuthTime:  jwt.NewNumericDate(data.AuthTime),
		Username:  data.Username,
		Name:      data.Name,
		Gender:    data.Gender,
		Birthdate: data.Birthdate.Format("02-01-2006"),
		Picture:   data.Picture,
	}

	token, err := createSignedTokenWithClaims(data.Key, claims, nil)
	if err != nil {
		return jwtclaims.IDTokenClaims{}, "", err
	}

	return claims, token, nil
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
