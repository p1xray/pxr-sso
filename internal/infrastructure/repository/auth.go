package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"log/slog"
)

type Auth struct {
	log     *slog.Logger
	storage storage.Storage
}

func NewAuthRepository(log *slog.Logger, storage storage.Storage) *Auth {
	return &Auth{
		log:     log,
		storage: storage,
	}
}

func (a *Auth) UserData(ctx context.Context, username string, clientCode string) (dto.User, error) {
	const op = "repository.UserDataByUsername"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	// Get user from storage.
	user, err := a.storage.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))
		}

		log.Error("failed to get user", sl.Err(err))

		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user permissions from storage.
	permissions, err := a.storage.PermissionsByUserID(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user permissions", sl.Err(err))

		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get user client with client code from storage.
	client, err := a.storage.UserClientByCode(ctx, user.ID, clientCode)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))
		}

		log.Error("failed to get user client", sl.Err(err))

		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	permissionCodes := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionCodes[i] = permission.Code
	}

	return dto.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Permissions:  permissionCodes,
		Client: dto.Client{
			ID:        client.ID,
			Code:      client.Code,
			SecretKey: client.SecretKey,
		},
	}, nil
}

// TODO: implement methods
// CreateSession(ctx context.Context, session entity.Session) error
