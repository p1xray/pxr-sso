package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/lib/token"
	clientcrud "pxr-sso/internal/logic/crud/client"
	sessioncrud "pxr-sso/internal/logic/crud/session"
	usercrud "pxr-sso/internal/logic/crud/user"
	"pxr-sso/internal/logic/dto"
	"pxr-sso/internal/logic/service"
	"pxr-sso/internal/storage"
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

// Login checks if user with given credentials exists in the system and returns access and refresh tokens.
func (a *Auth) Login(ctx context.Context, data dto.LoginDTO) (dto.TokensDTO, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
	)
	log.Info("attempting to login user")

	// Get user data from storage.
	user, err := a.userCRUD.UserWithPermissionByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		a.log.Warn("invalid credentials", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
	}

	// Get user client by user link and client code from storage.
	client, err := a.clientCRUD.UserClient(ctx, user.ID, data.ClientCode)
	if err != nil {
		a.log.Error("failed to get client", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access token.
	accessToken, err := token.NewAccessToken(&user, &client, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create refresh token.
	refreshToken, refreshTokenID, err := token.NewRefreshToken(client.SecretKey, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate refresh token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := dto.CreateSessionDTO{
		UserID:       user.ID,
		RefreshToken: refreshTokenID,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.sessionCRUD.CreateSession(ctx, sessionToCreate); err != nil {
		a.log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	return dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Register registers new user in the system and returns access and refresh tokens.
func (a *Auth) Register(ctx context.Context, data dto.RegisterDTO) (dto.TokensDTO, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
	)
	log.Info("attempting to register new user")

	// Get user from storage.
	user, err := a.userCRUD.UserWithPermissionByUsername(ctx, data.Username)
	if err != nil && !errors.Is(err, storage.ErrEntityNotFound) {
		a.log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check if user with given username already exists
	if user.ID > emptyValue {
		a.log.Warn("user with given username already exists", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserExists)
	}

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		a.log.Error("failed to get client", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Generate hash from password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create new user in storage.
	createUserData := dto.CreateUserDTO{
		Username:      data.Username,
		PasswordHash:  passwordHash,
		FIO:           data.FIO,
		DateOfBirth:   data.DateOfBirth,
		Gender:        data.Gender,
		AvatarFileKey: data.AvatarFileKey,
	}

	newUserID, err := a.userCRUD.CreateUser(ctx, createUserData)
	if err != nil {
		a.log.Error("failed to create user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create user's client link in storage.
	if err = a.userCRUD.CreateUserClientLink(ctx, newUserID, client.ID); err != nil {
		a.log.Error("failed to create user's client link", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user permissions from storage.
	newUser, err := a.userCRUD.UserWithPermission(ctx, newUserID)
	if err != nil {
		a.log.Error("failed to get user permissions", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access token.
	accessToken, err := token.NewAccessToken(&newUser, &client, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create refresh token.
	refreshToken, refreshTokenID, err := token.NewRefreshToken(client.SecretKey, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate refresh token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := dto.CreateSessionDTO{
		UserID:       newUser.ID,
		RefreshToken: refreshTokenID,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.sessionCRUD.CreateSession(ctx, sessionToCreate); err != nil {
		a.log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	return dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshTokens refreshes the user's auth tokens.
func (a *Auth) RefreshTokens(ctx context.Context, data dto.RefreshTokensDTO) (dto.TokensDTO, error) {
	const op = "auth.RefreshTokens"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
	)
	log.Info("attempting to refresh tokens")

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		a.log.Error("failed to get client", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := token.ParseRefreshToken(data.RefreshToken, client.SecretKey)
	if err != nil {
		a.log.Error("failed to parse refresh token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get session by refresh token ID from storage.
	session, err := a.sessionCRUD.SessionByRefreshToken(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			a.log.Warn("session not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrSessionNotFound)
		}

		a.log.Error("failed to get session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check session expiration time.
	now := time.Now()
	if session.ExpiresAt.Before(now) {
		a.log.Warn("refresh token expired")

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrRefreshTokenExpired)
	}

	// Check session user agent and fingerprint.
	if session.UserAgent != data.UserAgent && session.Fingerprint != data.Fingerprint {
		a.log.Warn("invalid session")

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrInvalidSession)
	}

	// Get user from storage.
	user, err := a.userCRUD.UserWithPermission(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Remove current session from storage.
	if err = a.sessionCRUD.RemoveSession(ctx, session.ID); err != nil {
		a.log.Error("failed to remove session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access token.
	accessToken, err := token.NewAccessToken(&user, &client, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create refresh token.
	refreshToken, refreshTokenID, err := token.NewRefreshToken(client.SecretKey, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate refresh token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := dto.CreateSessionDTO{
		UserID:       user.ID,
		RefreshToken: refreshTokenID,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.sessionCRUD.CreateSession(ctx, sessionToCreate); err != nil {
		a.log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("tokens refreshed successfully")

	return dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Logout terminates the user's session.
func (a *Auth) Logout(ctx context.Context, data dto.LogoutDTO) error {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
	)
	log.Info("attempting to user logout")

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		a.log.Error("failed to get client", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := token.ParseRefreshToken(data.RefreshToken, client.SecretKey)
	if err != nil {
		a.log.Error("failed to parse refresh token", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Get session by refresh token from storage.
	session, err := a.sessionCRUD.SessionByRefreshToken(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			a.log.Warn("session not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, service.ErrSessionNotFound)
		}

		a.log.Error("failed to get session", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Remove current session from storage.
	if err = a.sessionCRUD.RemoveSession(ctx, session.ID); err != nil {
		a.log.Error("failed to remove session", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logout successfully")

	return nil
}
