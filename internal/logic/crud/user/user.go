package usercrud

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/logic/crud"
	"pxr-sso/internal/logic/dto"
	"pxr-sso/internal/storage"
)

// CRUD provides methods for managing user data.
type CRUD struct {
	log                *slog.Logger
	userProvider       crud.UserProvider
	userSaver          crud.UserSaver
	permissionProvider crud.PermissionProvider
}

// New creates a new instance of the user's CRUD.
func New(
	log *slog.Logger,
	userProvider crud.UserProvider,
	userSaver crud.UserSaver,
	permissionProvider crud.PermissionProvider,
) *CRUD {
	return &CRUD{
		log:                log,
		userProvider:       userProvider,
		userSaver:          userSaver,
		permissionProvider: permissionProvider,
	}
}

// UserWithPermission returns user data with permissions by user ID.
func (c *CRUD) UserWithPermission(ctx context.Context, id int64) (dto.UserDTO, error) {
	const op = "usercrud.UserWithPermission"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("user ID", id),
	)

	user, err := c.userProvider.User(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found in storage", sl.Err(err))

			return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get user from storage", sl.Err(err))

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	permissionCodes, err := c.userPermissionCodes(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user permissions from storage", sl.Err(err))

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userData := dto.UserDTO{
		ID:          user.ID,
		Permissions: permissionCodes,
	}

	return userData, nil
}

// UserWithPermissionByUsername returns user data with permissions by username.
func (c *CRUD) UserWithPermissionByUsername(ctx context.Context, username string) (dto.UserDTO, error) {
	const op = "usercrud.UserWithPermissionByUsername"

	log := c.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	user, err := c.userProvider.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found in storage", sl.Err(err))

			return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get user from storage", sl.Err(err))

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	permissionCodes, err := c.userPermissionCodes(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user permissions from storage", sl.Err(err))

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userData := dto.UserDTO{
		ID:          user.ID,
		Permissions: permissionCodes,
	}

	return userData, nil
}

func (c *CRUD) userPermissionCodes(ctx context.Context, userID int64) ([]string, error) {
	userPermissions, err := c.permissionProvider.UserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var permissionCodes []string
	for _, permission := range userPermissions {
		permissionCodes = append(permissionCodes, permission.Code)
	}

	return permissionCodes, nil
}
