package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/logic/service"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const emptyValue = 0

type Auth struct {
	Session Session

	user   User
	client dto.Client

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuth(user dto.User, accessTokenTTL, refreshTokenTTL time.Duration) Auth {
	return Auth{
		user:            NewUser(user.ID, user.Username, user.PasswordHash, user.Permissions),
		client:          user.Client,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (a *Auth) Tokens() Tokens {
	return a.Session.Tokens
}

func (a *Auth) Login(data LoginParams) error {
	const op = "entity.Auth.Login"

	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(a.user.PasswordHash), []byte(data.Password)); err != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Create new session.
	createNewSessionParams := CreateNewSessionParams{
		UserID:          a.user.ID,
		UserPermissions: a.user.Permissions,
		ClientCode:      a.client.Code,
		ClientSecretKey: a.client.SecretKey,
		Issuer:          data.Issuer,
		UserAgent:       data.UserAgent,
		Fingerprint:     data.Fingerprint,
		AccessTokenTTL:  a.accessTokenTTL,
		RefreshTokenTTL: a.refreshTokenTTL,
	}
	session, err := NewSession(createNewSessionParams)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.Session = session

	return nil
}

func (a *Auth) Register(data RegisterParams) error {
	const op = "entity.Auth.Register"

	// Check if user with given username already exists.
	if a.user.ID > emptyValue {
		return fmt.Errorf("%s: %w", op, ErrUserExists)
	}

	// Generate hash from password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrGeneratePasswordHash)
	}
}
