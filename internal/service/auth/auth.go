package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/guregu/null/v6"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"pxr-sso/internal/domain"
	"pxr-sso/internal/dto"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/lib/token"
	"pxr-sso/internal/service"
	"pxr-sso/internal/storage"
	"time"
)

const (
	emptyValue = 0
)

// Auth is service for working with user authentication and authorization.
type Auth struct {
	log             *slog.Logger
	storage         storage.SSOStorage
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// New creates a new auth service.
func New(
	log *slog.Logger,
	storage storage.SSOStorage,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:             log,
		storage:         storage,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
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

	// Get user from storage.
	user, err := a.storage.UserByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
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
	client, err := a.storage.UserClient(ctx, user.ID, data.ClientCode)
	if err != nil {
		a.log.Error("failed to get client", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user permissions from storage.
	userPermissions, err := a.storage.UserPermissions(ctx, user.ID)
	if err != nil {
		a.log.Error("failed to get user permissions", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access token.
	var permissionCodes []string
	for _, permission := range userPermissions {
		permissionCodes = append(permissionCodes, permission.Code)
	}

	userData := &dto.UserDTO{
		ID:          user.ID,
		Permissions: permissionCodes,
	}

	clientData := &dto.ClientDTO{
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}

	accessToken, err := token.NewAccessToken(userData, clientData, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create refresh token.
	refreshToken := token.NewRefreshToken()

	// Create session in storage.
	sessionToCreate := domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.storage.CreateSession(ctx, sessionToCreate); err != nil {
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
	user, err := a.storage.UserByUsername(ctx, data.Username)
	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		a.log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check if user with given username already exists
	if user.ID > emptyValue {
		a.log.Warn("user with given username already exists", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserExists)
	}

	// Get client by code from storage.
	client, err := a.storage.ClientByCode(ctx, data.ClientCode)
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
	userToCreate := domain.User{
		Username:      data.Username,
		PasswordHash:  string(passwordHash),
		FIO:           data.FIO,
		DateOfBirth:   null.TimeFromPtr(data.DateOfBirth),
		Gender:        data.Gender.ToNullInt16(),
		AvatarFileKey: null.StringFromPtr(data.AvatarFileKey),
	}

	newUserID, err := a.storage.CreateUser(ctx, userToCreate)
	if err != nil {
		a.log.Error("failed to create user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create user's client link in storage.
	if err = a.storage.CreateUserClientLink(ctx, newUserID, client.ID); err != nil {
		a.log.Error("failed to create user's client link", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user permissions from storage.
	userPermissions, err := a.storage.UserPermissions(ctx, newUserID)
	if err != nil {
		a.log.Error("failed to get user permissions", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access token.
	var permissionCodes []string
	for _, permission := range userPermissions {
		permissionCodes = append(permissionCodes, permission.Code)
	}

	userData := &dto.UserDTO{
		ID:          newUserID,
		Permissions: permissionCodes,
	}

	clientData := &dto.ClientDTO{
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}

	accessToken, err := token.NewAccessToken(userData, clientData, a.accessTokenTTL, data.Issuer)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create refresh token.
	refreshToken := token.NewRefreshToken()

	// Create session in storage.
	sessionToCreate := domain.Session{
		UserID:       newUserID,
		RefreshToken: refreshToken,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.storage.CreateSession(ctx, sessionToCreate); err != nil {
		a.log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	return dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
