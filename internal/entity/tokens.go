package entity

import (
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"strconv"
)

type Tokens struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenID string
}

func NewTokens(data CreateTokensParams) (Tokens, error) {
	// Create access token.
	createAccessTokenData := jwtcreator.AccessTokenCreateData{
		Subject:  strconv.FormatInt(data.UserID, 10),
		Audience: data.ClientCode,
		Scopes:   data.Permissions,
		Issuer:   data.Issuer,
		TTL:      data.AccessTokenTTL,
		Key:      []byte(data.SecretKey),
	}
	accessToken, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		return Tokens{}, err
	}

	// Create refresh token.
	refreshToken, refreshTokenID, err := jwtcreator.NewRefreshToken([]byte(data.SecretKey), data.RefreshTokenTTL)
	if err != nil {

		return Tokens{}, err
	}

	return Tokens{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenID: refreshTokenID,
	}, nil
}
