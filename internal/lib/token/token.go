package token

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"pxr-sso/internal/logic/dto"
	jwtmiddleware "pxr-sso/pkg/jwt"
	"strconv"
	"strings"
	"time"
)

// NewAccessToken returns new JWT with claims.
func NewAccessToken(user *dto.UserDTO, client *dto.ClientDTO, ttl time.Duration, issuer string) (string, error) {
	now := time.Now()
	claims := jwtmiddleware.AccessTokenClaims{
		jwt.Claims{
			ID:        uuid.New().String(),
			Subject:   strconv.FormatInt(user.ID, 10),
			Issuer:    issuer,
			Audience:  []string{client.Code},
			Expiry:    jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		jwtmiddleware.RegisteredCustomClaims{
			Scope: strings.Join(user.Permissions, " "),
		},
	}

	tokenStr, err := createSignedTokenWithClaims(claims, client.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// NewRefreshToken returns new refresh token.
func NewRefreshToken(secretKey string, ttl time.Duration) (refreshToken string, refreshTokenID string, err error) {
	id := uuid.New().String()
	now := time.Now()
	claims := jwtmiddleware.RefreshTokenClaims{
		ID:     id,
		Expiry: jwt.NewNumericDate(now.Add(ttl)),
	}

	tokenStr, err := createSignedTokenWithClaims(claims, secretKey)
	if err != nil {
		return "", "", err
	}

	return tokenStr, id, nil
}

// ParseAccessToken parses access token as a string using a secret key into a set of claims.
func ParseAccessToken(tokenStr string, secretKey string) (jwtmiddleware.AccessTokenClaims, error) {
	token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return jwtmiddleware.AccessTokenClaims{}, err
	}

	claims := jwtmiddleware.AccessTokenClaims{}
	if err = token.Claims(secretKey, &claims); err != nil {
		return jwtmiddleware.AccessTokenClaims{}, err
	}

	return claims, nil
}

// ParseRefreshToken parses refresh token as a string using a secret key into a set of claims.
func ParseRefreshToken(tokenStr string, secretKey string) (jwtmiddleware.RefreshTokenClaims, error) {
	token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return jwtmiddleware.RefreshTokenClaims{}, err
	}

	claims := jwtmiddleware.RefreshTokenClaims{}
	if err = token.Claims(secretKey, &claims); err != nil {
		return jwtmiddleware.RefreshTokenClaims{}, err
	}

	return claims, nil
}

func createSignedTokenWithClaims(claims interface{}, secretKey string) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte(secretKey)},
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
