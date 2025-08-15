package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	Sessions []Session
	User     User

	client          dto.Client
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuth(accessTokenTTL, refreshTokenTTL time.Duration, setters ...AuthOption) (Auth, error) {
	auth := Auth{
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}

	for _, setter := range setters {
		if err := setter(&auth); err != nil {
			return Auth{}, err
		}
	}

	return auth, nil
}

func (a *Auth) Login(data LoginParams) (Tokens, error) {
	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(a.User.PasswordHash), []byte(data.Password)); err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
	}

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateSession, err)
	}

	return tokens, nil
}

func (a *Auth) Register(data RegisterParams) (Tokens, error) {
	// Check if user with given username already exists.
	if a.User.ID > emptyID {
		return Tokens{}, ErrUserExists
	}

	// Generate hash from password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrGeneratePasswordHash, err)
	}

	// Create new user.
	user := NewUser(
		data.Username,
		data.FullName,
		data.DateOfBirth,
		data.Gender,
		data.AvatarFileKey,
		WithUserPasswordHash(string(passwordHash)),
		WithUserRoles([]string{defaultRoleName}), // TODO: get this from storage.
		WithUserPermissions([]string{}),          // TODO: get permissions from storage by default role.
	)
	user.SetToCreate()
	a.setUser(user)

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateSession, err)
	}

	return tokens, nil
}

func (a *Auth) RefreshTokens(data RefreshTokensParams) (Tokens, error) {
	// Check session data.
	for _, session := range a.Sessions {
		if err := session.Validate(data.UserAgent, data.Fingerprint); err != nil {
			return Tokens{}, fmt.Errorf("%w: %w", ErrValidateSession, err)
		}
	}

	// Set current session to remove.
	for i := range a.Sessions {
		a.Sessions[i].SetToRemove()
	}

	// Create new session.
	tokens, err := a.createNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateSession, err)
	}

	return tokens, nil
}

func (a *Auth) Logout() error {
	// Check if session exist.
	if len(a.Sessions) == 0 {
		return ErrSessionNotFound
	}

	// Set current session to remove.
	for i := range a.Sessions {
		a.Sessions[i].SetToRemove()
	}

	return nil
}

func (a *Auth) createNewSession(issuer, userAgent, fingerprint string) (Tokens, error) {
	generateTokensParams := SessionWithGeneratedTokensParams{
		UserPermissions: a.User.Permissions,
		ClientCode:      a.client.Code,
		ClientSecretKey: a.client.SecretKey,
		Issuer:          issuer,
		AccessTokenTTL:  a.accessTokenTTL,
		RefreshTokenTTL: a.refreshTokenTTL,
	}

	session, err := NewSession(
		a.User.ID,
		userAgent,
		fingerprint,
		WithGeneratedTokens(generateTokensParams),
	)
	if err != nil {
		return Tokens{}, err
	}

	session.SetToCreate()

	a.addSession(session)

	return session.Tokens, nil
}

func (a *Auth) addSession(session Session) {
	a.Sessions = append(a.Sessions, session)
}

func (a *Auth) setUser(user User) {
	a.User = user
}
