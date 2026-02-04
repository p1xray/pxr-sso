package entity

import (
	"fmt"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"strconv"
)

// Tokens is the user session tokens entity.
type Tokens struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenID string
}

// NewTokens returns new user session tokens entity.
func NewTokens(data CreateTokensParams) (Tokens, error) {
	// Create access token.
	createAccessTokenData := jwtcreator.AccessTokenCreateData{
		Subject:   strconv.FormatInt(data.UserID, 10),
		Audiences: data.Audiences,
		Scopes:    data.Permissions,
		Issuer:    data.Issuer,
		TTL:       data.AccessTokenTTL,
		Key:       []byte(data.SecretKey),
	}
	_, accessToken, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateAccessToken, err)
	}

	// Create refresh token.
	refreshTokenClaims, refreshToken, err := jwtcreator.NewRefreshToken([]byte(data.SecretKey), data.RefreshTokenTTL)
	if err != nil {

		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateRefreshToken, err)
	}

	return Tokens{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenID: refreshTokenClaims.ID,
	}, nil
}
