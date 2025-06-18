package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	"github.com/p1xray/pxr-sso/internal/logic/service"
	"github.com/p1xray/pxr-sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

// Register registers new user in the system and returns access and refresh tokens.
func (a *Auth) Register(ctx context.Context, data dto.RegisterDTO) (dto.TokensDTO, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
	)
	log.Info("attempting to register new user")

	// Get user from storage.
	user, err := a.userCRUD.UserWithPermissionsByUsername(ctx, data.Username)
	if err != nil && !errors.Is(err, storage.ErrEntityNotFound) {
		log.Error("failed to get user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Check if user with given username already exists
	if user.ID > emptyValue {
		log.Warn("user with given username already exists", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserExists)
	}

	// Get client by code from storage.
	client, err := a.clientCRUD.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrClientNotFound)
		}

		log.Error("failed to get client", sl.Err(err))

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
		log.Error("failed to create user", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create user's client link in storage.
	if err = a.userCRUD.CreateUserClientLink(ctx, newUserID, client.ID); err != nil {
		log.Error("failed to create user's client link", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user permissions from storage.
	newUser, err := a.userCRUD.UserWithPermissions(ctx, newUserID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("failed to get user permissions", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create access and refresh tokens.
	tokensData, err := a.createAccessAndRefreshTokens(log, user, client, data.Issuer)
	if err != nil {
		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create session in storage.
	sessionToCreate := dto.CreateSessionDTO{
		UserID:       newUser.ID,
		RefreshToken: tokensData.RefreshTokenID,
		UserAgent:    data.UserAgent,
		Fingerprint:  data.Fingerprint,
		ExpiresAt:    time.Now().Add(a.refreshTokenTTL),
	}
	if err = a.sessionCRUD.CreateSession(ctx, sessionToCreate); err != nil {
		log.Error("failed to create session", sl.Err(err))

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	return dto.TokensDTO{AccessToken: tokensData.AccessToken, RefreshToken: tokensData.RefreshToken}, nil
}
