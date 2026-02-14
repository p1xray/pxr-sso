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
	issuer          string
}

func NewToken(cfg TokenConfig) *Token {
	return &Token{
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (t *Token) GenerateTokens(
	scope []string,
	audience string,
	user dto.User,
	client dto.Client,
) (dto.Token, error) {
	const op = "generate tokens"

	accessTokenClaims, accessToken, err := t.generateAccessToken(scope, audience, user, client)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, refreshToken, err := t.generateRefreshToken(client.SecretKey())
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, idToken, err := t.generateIDToken(scope)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
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

func (t *Token) generateIDToken(scope []string) (jwtclaims.RefreshTokenClaims, string, error) {
	const op = "id token"

	// TODO: implement this
	return jwtclaims.RefreshTokenClaims{}, "", nil
}
