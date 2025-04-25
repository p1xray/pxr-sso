package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"pxr-sso/internal/dto"
	"strconv"
	"time"
)

// CustomClaims are custom claims of the current SSO project.
type CustomClaims struct {
	ClientID    string   `json:"client_id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// NewAccessToken returns new JWT with claims.
func NewAccessToken(user *dto.User, client *dto.Client, ttl time.Duration, issuer string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		client.Code,
		user.Permissions,
		jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   strconv.FormatInt(user.ID, 10),
			Issuer:    issuer,
			Audience:  []string{client.Code},
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(client.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// NewRefreshToken returns new refresh token.
func NewRefreshToken() string {
	token := uuid.New()
	return token.String()
}
