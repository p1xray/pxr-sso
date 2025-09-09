package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// maxUserSessionsCount specifies how many sessions a user can have at the same time.
// If the number of sessions exceeds this value, all previous user sessions are removed from the storage.
const maxUserSessionsCount int = 5

// Auth is the user authentication entity.
type Auth struct {
	Sessions []Session
	User     User

	client                 dto.Client
	defaultRoles           []dto.Role
	defaultPermissionCodes []string
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
}

// NewAuth returns a new auth entity.
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

// Login verifies the user's login data, and if successful, creates a new user session.
func (a *Auth) Login(data LoginParams) (Tokens, error) {
	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(a.User.PasswordHash), []byte(data.Password)); err != nil {
		return Tokens{}, fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
	}

	// Check user sessions count.
	if len(a.Sessions) >= maxUserSessionsCount {
		// Set all sessions to remove.
		for i := range a.Sessions {
			a.Sessions[i].SetToRemove()
		}
	}

	// Create new session.
	tokens, err := a.CreateNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

// Register creates a new user in the system.
func (a *Auth) Register(data RegisterParams) error {
	// Check if user with given username already exists.
	if a.User.ID > emptyID {
		return ErrUserExists
	}

	// Generate hash from password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGeneratePasswordHash, err)
	}

	// Create new user.
	user := NewUser(
		data.Username,
		data.FullName,
		data.DateOfBirth,
		data.Gender,
		data.AvatarFileKey,
		WithUserPasswordHash(string(passwordHash)),
		WithUserRoles(a.defaultRoles),
		WithUserPermissions(a.defaultPermissionCodes),
	)
	user.SetToCreate()
	a.setUser(user)

	return nil
}

// RefreshTokens refreshes the user's tokens, and if successful creates a new user session.
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
	tokens, err := a.CreateNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

// Logout deletes the current user session.
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

// CreateNewSession creates a new user session.
func (a *Auth) CreateNewSession(issuer, userAgent, fingerprint string) (Tokens, error) {
	generateTokensParams := SessionWithGeneratedTokensParams{
		UserPermissions: a.User.Permissions,
		Audiences:       a.client.Audiences,
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
		return Tokens{}, fmt.Errorf("%w: %w", ErrCreateSession, err)
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

func (a *Auth) ClientID() int64 {
	return a.client.ID
}

func (a *Auth) ClientCode() string {
	return a.client.Code
}
