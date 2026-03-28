package generator

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"slices"
	"time"
)

type Token struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	idTokenTTL      time.Duration
	issuer          string
}

func NewToken(cfg TokenConfig) *Token {
	return &Token{
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
		idTokenTTL:      cfg.IDTokenTTL,
		issuer:          cfg.Issuer,
	}
}

func (t *Token) GenerateTokens(
	scope []string,
	audience string,
	user dto.User,
	client dto.Client,
) (dto.Token, error) {
	const op = "tokens"

	accessTokenClaims, accessToken, err := t.generateAccessToken(scope, audience, user, client)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	_, refreshToken, err := t.generateRefreshToken(client.SecretKey())
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	_, idToken, err := t.generateIDToken(user, client)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	tokens := dto.NewToken(
		accessToken,
		accessTokenClaims.TokenType,
		refreshToken,
		idToken,
		jwtclaims.NumericDateToInt64(accessTokenClaims.Expiry),
	)

	return tokens, nil
}

func (t *Token) generateAccessToken(
	scope []string,
	audience string,
	user dto.User,
	client dto.Client,
) (jwtclaims.AccessTokenClaims, string, error) {
	const op = "access token"

	permissions := make([]string, 0, len(user.Roles()))
	for _, role := range user.Roles() {
		slices.Concat(permissions, role.Permissions())
	}
	scopes := extslices.Union(scope, permissions)

	createAccessTokenData := jwtcreator.AccessTokenCreateData{
		Subject:   user.IDString(),
		Audiences: []string{audience},
		Scopes:    scopes,
		Issuer:    t.issuer,
		TTL:       t.accessTokenTTL,
		Key:       []byte(client.SecretKey()),
	}
	claims, accessToken, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return claims, accessToken, nil
}

func (t *Token) generateRefreshToken(key string) (jwtclaims.RefreshTokenClaims, string, error) {
	const op = "refresh token"

	claims, refreshToken, err := jwtcreator.NewRefreshToken([]byte(key), t.refreshTokenTTL)
	if err != nil {
		return jwtclaims.RefreshTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return claims, refreshToken, nil
}

func (t *Token) generateIDToken(user dto.User, client dto.Client) (jwtclaims.IDTokenClaims, string, error) {
	const op = "id token"

	createIDTokenData := jwtcreator.IDTokenCreateData{
		Subject:   user.IDString(),
		ClientID:  client.Code(),
		Issuer:    t.issuer,
		AuthTime:  time.Now(),
		Username:  user.Username(),
		Name:      user.FullName(),
		Gender:    "",                                                       // TODO: implement user.Gender(),
		Birthdate: time.Date(2000, time.November, 12, 0, 0, 0, 0, time.UTC), // TODO: implement user.Birthdate(),
		Picture:   "",                                                       // TODO: implement user.PictureURL(),
		TTL:       t.idTokenTTL,
		Key:       []byte(client.SecretKey()),
	}
	claims, idToken, err := jwtcreator.NewIDToken(createIDTokenData)
	if err != nil {
		return jwtclaims.IDTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return claims, idToken, nil
}
