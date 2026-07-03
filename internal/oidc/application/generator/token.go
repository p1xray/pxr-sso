package generator

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"time"
)

type TokensGenerator interface {
	Generate(scope []string, audience string, user dto.User, client dto.Client) (dto.Tokens, error)
}

type tokens struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	idTokenTTL      time.Duration
	issuer          string
}

func NewTokensGenerator(cfg TokenConfig) *tokens {
	return &tokens{
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
		idTokenTTL:      cfg.IDTokenTTL,
		issuer:          cfg.Issuer,
	}
}

func (t *tokens) Generate(
	scope []string,
	audience string,
	user dto.User,
	client dto.Client,
) (dto.Tokens, error) {
	const op = "tokens"

	accessToken, err := t.generateAccessToken(scope, audience, user, client)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	refreshToken, err := t.generateRefreshToken(client.SecretKey())
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	idToken, err := t.generateIDToken(user, client)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	generatedTokens := dto.NewTokens(
		accessToken,
		refreshToken,
		idToken,
	)

	return generatedTokens, nil
}

func (t *tokens) generateAccessToken(
	scope []string,
	audience string,
	user dto.User,
	client dto.Client,
) (dto.Token, error) {
	const op = "access token"

	permissions := make([]string, 0)
	for _, role := range user.Roles() {
		for _, permission := range role.Permissions() {
			permissions = append(permissions, permission.Code())
		}
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
	claims, accessTokenString, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		claims.ID,
		claims.TokenType,
		accessTokenString,
		jwtclaims.NumericDateToInt64(claims.Expiry),
	)

	return token, nil
}

func (t *tokens) generateRefreshToken(key string) (dto.Token, error) {
	const op = "refresh token"

	claims, refreshTokenString, err := jwtcreator.NewRefreshToken([]byte(key), t.refreshTokenTTL)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		claims.ID,
		claims.TokenType,
		refreshTokenString,
		jwtclaims.NumericDateToInt64(claims.Expiry),
	)

	return token, nil
}

func (t *tokens) generateIDToken(user dto.User, client dto.Client) (dto.Token, error) {
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
	claims, idTokenString, err := jwtcreator.NewIDToken(createIDTokenData)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		claims.ID,
		claims.TokenType,
		idTokenString,
		jwtclaims.NumericDateToInt64(claims.Expiry),
	)

	return token, nil
}
