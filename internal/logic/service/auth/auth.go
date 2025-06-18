package auth

import (
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	clientcrud "github.com/p1xray/pxr-sso/internal/logic/crud/client"
	sessioncrud "github.com/p1xray/pxr-sso/internal/logic/crud/session"
	usercrud "github.com/p1xray/pxr-sso/internal/logic/crud/user"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"log/slog"
	"strconv"
	"time"
)

const (
	emptyValue = 0
)

// Auth is service for working with user authentication and authorization.
type Auth struct {
	log             *slog.Logger
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	userCRUD        *usercrud.CRUD
	clientCRUD      *clientcrud.CRUD
	sessionCRUD     *sessioncrud.CRUD
}

// New creates a new auth service.
func New(
	log *slog.Logger,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	userCRUD *usercrud.CRUD,
	clientCRUD *clientcrud.CRUD,
	sessionCRUD *sessioncrud.CRUD,
) *Auth {
	return &Auth{
		log:             log,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		userCRUD:        userCRUD,
		clientCRUD:      clientCRUD,
		sessionCRUD:     sessionCRUD,
	}
}

type tokens struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenID string
}

func (a *Auth) createAccessAndRefreshTokens(
	log *slog.Logger,
	user dto.UserWithPermissionsDTO,
	client dto.ClientDTO,
	issuer string,
) (tokens, error) {
	// Create access token.
	createAccessTokenData := jwtcreator.AccessTokenCreateData{
		Subject:  strconv.FormatInt(user.ID, 10),
		Audience: client.Code,
		Scopes:   user.Permissions,
		Issuer:   issuer,
		TTL:      a.accessTokenTTL,
		Key:      []byte(client.SecretKey),
	}
	accessToken, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))

		return tokens{}, err
	}

	// Create refresh token.
	refreshToken, refreshTokenID, err := jwtcreator.NewRefreshToken([]byte(client.SecretKey), a.refreshTokenTTL)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))

		return tokens{}, err
	}

	tokensData := tokens{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenID: refreshTokenID,
	}

	return tokensData, nil
}
