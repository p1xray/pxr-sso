package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	emptyValue      = 0
	defaultRoleName = "member"
)

type Auth struct {
	Sessions []Session
	User     User

	client          dto.Client
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuth(accessTokenTTL, refreshTokenTTL time.Duration, setters ...AuthOption) Auth {
	auth := Auth{
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}

	for _, setter := range setters {
		setter(&auth)
	}

	return auth
}

func (a *Auth) Login(data LoginParams) (Tokens, error) {
	const op = "entity.Auth.Login"

	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(a.User.PasswordHash), []byte(data.Password)); err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *Auth) Register(data RegisterParams) (Tokens, error) {
	const op = "entity.Auth.Register"

	// Check if user with given username already exists.
	if a.User.ID > emptyValue {
		return Tokens{}, fmt.Errorf("%s: %w", op, ErrUserExists)
	}

	// Generate hash from password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, ErrGeneratePasswordHash)
	}

	// Create new user.
	user := NewUser(
		0,
		data.Username,
		string(passwordHash),
		data.FullName,
		data.DateOfBirth,
		data.Gender,
		data.AvatarFileKey,
		[]string{defaultRoleName}, // TODO: get this from storage.
		[]string{},                // TODO: get permissions from storage by default role.
	)
	user.SetToCreate()
	a.setUser(user)

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *Auth) RefreshTokens(data RefreshTokensParams) (Tokens, error) {
	const op = "entity.Auth.RefreshTokens"

	// Check session data.
	for _, session := range a.Sessions {
		if err := session.Validate(data.UserAgent, data.Fingerprint); err != nil {
			return Tokens{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	// Set current session to remove.
	for _, session := range a.Sessions {
		session.SetToRemove()
	}

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *Auth) Logout() error {
	const op = "entity.Auth.Logout"

	// Check if session exist.
	if len(a.Sessions) == 0 {
		return fmt.Errorf("%s: %w", op, ErrSessionNotFound)
	}

	// Set current session to remove.
	for _, session := range a.Sessions {
		session.SetToRemove()
	}

	return nil
}

func (a *Auth) createNewSession(issuer, userAgent, fingerprint string) (Tokens, error) {
	createNewSessionParams := CreateNewSessionParams{
		UserID:          a.User.ID,
		UserPermissions: a.User.Permissions,
		ClientCode:      a.client.Code,
		ClientSecretKey: a.client.SecretKey,
		Issuer:          issuer,
		UserAgent:       userAgent,
		Fingerprint:     fingerprint,
		AccessTokenTTL:  a.accessTokenTTL,
		RefreshTokenTTL: a.refreshTokenTTL,
	}
	session, err := NewSession(createNewSessionParams)
	if err != nil {
		return Tokens{}, err
	}

	a.addSession(session)

	return session.Tokens, nil
}

func (a *Auth) addSession(session Session) {
	a.Sessions = append(a.Sessions, session)
}

func (a *Auth) setUser(user User) {
	a.User = user
}
