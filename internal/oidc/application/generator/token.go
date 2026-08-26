package generator

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"slices"
	"time"
)

const (
	tokenTypeBearer  = "Bearer"
	tokenTypeRefresh = "Refresh"
	tokenTypeID      = "OpenID"
)

type TokensGenerator interface {
	Generate(
		scope []string,
		audience string,
		user dto.User,
		client dto.Client,
		session dto.AuthorizedSession,
	) (dto.Tokens, error)
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
	session dto.AuthorizedSession,
) (dto.Tokens, error) {
	const op = "tokens"

	accessToken, err := t.generateAccessToken(scope, audience, user, client, session)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	refreshToken := dto.Token{}
	if slices.Contains(scope, "offline_access") {
		refreshToken, err = t.generateRefreshToken(scope, audience, user, client, session)
		if err != nil {
			return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
	}

	idToken := dto.Token{}
	if slices.Contains(scope, "openid") {
		idToken, err = t.generateIDToken(scope, audience, user, client, session)
		if err != nil {
			return dto.Tokens{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
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
	session dto.AuthorizedSession,
) (dto.Token, error) {
	const op = "access token"

	permissions := make([]string, 0)
	for _, role := range user.Roles() {
		for _, permission := range role.Permissions() {
			permissions = append(permissions, permission.Code())
		}
	}
	scopes := extslices.Union(scope, permissions)

	tokenClaims := claims.NewAccessTokenClaims(
		t.issuer,
		user.IDString(),
		[]string{audience},
		t.accessTokenTTL,
		claims.WithAccessTokenAuthentication(session.AuthTime()),
		claims.WithAccessTokenAuthorization(scopes),
		claims.WithAccessTokenClient(client.Code()),
	)

	rawToken, err := jwt.CreateAccessToken(tokenClaims, []byte(client.SecretKey()))
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		tokenClaims.ID,
		tokenTypeBearer,
		rawToken,
		claims.NumericDateToInt64(tokenClaims.Expiry),
	)

	return token, nil
}

func (t *tokens) generateRefreshToken(
	scopes []string,
	audience string,
	user dto.User,
	client dto.Client,
	session dto.AuthorizedSession,
) (dto.Token, error) {
	const op = "refresh token"

	tokenClaims := claims.NewRefreshTokenClaims(
		t.issuer,
		user.IDString(),
		[]string{audience},
		t.refreshTokenTTL,
		claims.WithRefreshTokenAuthentication(session.AuthTime()),
		claims.WithRefreshTokenAuthorization(scopes),
		claims.WithRefreshTokenClient(client.Code()),
		claims.WithRefreshTokenSession(session.ID()),
	)

	rawToken, err := jwt.CreateRefreshToken(tokenClaims, []byte(client.SecretKey()))
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		tokenClaims.ID,
		tokenTypeRefresh,
		rawToken,
		claims.NumericDateToInt64(tokenClaims.Expiry),
	)

	return token, nil
}

func (t *tokens) generateIDToken(
	scopes []string,
	audience string,
	user dto.User,
	client dto.Client,
	session dto.AuthorizedSession,
) (dto.Token, error) {
	const op = "id token"

	claimsOptions := []claims.IDTokenClaimsOption{
		claims.WithIDTokenAuthentication(session.AuthTime()),
	}

	if slices.Contains(scopes, "profile") {
		// TODO: add this fields into db tables
		const (
			name              = "John Michael Doe"
			familyName        = "Doe"
			givenName         = "John"
			middleName        = "Michael"
			nickname          = "johndoe"
			preferredUsername = "Johnny"
			profile           = "https://example.com/profile/johndoe"
			picture           = "https://example.com/profile/picture/johndoe"
			website           = "https://johndoe.dev"
			gender            = "male"
			birthdate         = "1990-01-01"
			zoneInfo          = "America/New_York"
			locale            = "en-US"
		)

		updatedAt := time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC)

		profileClaimsOption := claims.WithIDTokenProfile(
			claims.WithProfileName(name, familyName, givenName, middleName),
			claims.WithProfileNameAlias(nickname, preferredUsername),
			claims.WithProfileURL(profile, picture, website),
			claims.WithProfileGender(gender),
			claims.WithProfileBirthdate(birthdate),
			claims.WithProfileLocation(zoneInfo, locale),
			claims.WithProfileUpdatedAt(updatedAt),
		)

		claimsOptions = append(claimsOptions, profileClaimsOption)
	}

	if slices.Contains(scopes, "email") {
		// TODO: add this fields into db tables
		const (
			email         = "johndoe@example.com"
			emailVerified = false
		)

		emailClaimsOption := claims.WithIDTokenEmail(email, emailVerified)

		claimsOptions = append(claimsOptions, emailClaimsOption)
	}

	if slices.Contains(scopes, "phone") {
		// TODO: add this fields into db tables
		const (
			phoneNumber         = "+1234567890"
			phoneNumberVerified = true
		)

		phoneClaimsOption := claims.WithIDTokenPhoneNumber(phoneNumber, phoneNumberVerified)

		claimsOptions = append(claimsOptions, phoneClaimsOption)
	}

	if slices.Contains(scopes, "address") {
		// TODO: add this fields into db tables
		const (
			addressFormatted = "123 Main St\nMetropolis, NY 10001\nUSA"
			streetAddress    = "123 Main St"
			locality         = "Metropolis"
			region           = "NY"
			postalCode       = "10001"
			country          = "USA"
		)

		addressClaimsOption := claims.WithIDTokenAddress(
			addressFormatted,
			streetAddress,
			locality,
			region,
			postalCode,
			country,
		)

		claimsOptions = append(claimsOptions, addressClaimsOption)
	}

	tokenClaims := claims.NewIDTokenClaims(
		t.issuer,
		user.IDString(),
		[]string{audience},
		t.refreshTokenTTL,
		claimsOptions...,
	)

	rawToken, err := jwt.CreateIDToken(tokenClaims, []byte(client.SecretKey()))
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(
		tokenClaims.ID,
		tokenTypeID,
		rawToken,
		claims.NumericDateToInt64(tokenClaims.Expiry),
	)

	return token, nil
}
