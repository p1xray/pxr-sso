package usercrud

import (
	"context"
	"github.com/guregu/null/v6"
	"pxr-sso/internal/domain"
	"pxr-sso/internal/logic/crud"
	"pxr-sso/internal/logic/dto"
	"time"
)

// CRUD provides methods for managing user data.
type CRUD struct {
	userProvider       crud.UserProvider
	userSaver          crud.UserSaver
	permissionProvider crud.PermissionProvider
}

// New creates a new instance of the user's CRUD.
func New(
	userProvider crud.UserProvider,
	userSaver crud.UserSaver,
	permissionProvider crud.PermissionProvider,
) *CRUD {
	return &CRUD{
		userProvider:       userProvider,
		userSaver:          userSaver,
		permissionProvider: permissionProvider,
	}
}

// UserWithPermission returns user data with permissions by user ID.
func (c *CRUD) UserWithPermission(ctx context.Context, id int64) (dto.UserDTO, error) {
	user, err := c.userProvider.User(ctx, id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	permissionCodes, err := c.userPermissionCodes(ctx, user.ID)
	if err != nil {
		return dto.UserDTO{}, err
	}

	userData := dto.UserDTO{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Permissions:  permissionCodes,
	}

	return userData, nil
}

// UserWithPermissionByUsername returns user data with permissions by username.
func (c *CRUD) UserWithPermissionByUsername(ctx context.Context, username string) (dto.UserDTO, error) {
	user, err := c.userProvider.UserByUsername(ctx, username)
	if err != nil {
		return dto.UserDTO{}, err
	}

	permissionCodes, err := c.userPermissionCodes(ctx, user.ID)
	if err != nil {
		return dto.UserDTO{}, err
	}

	userData := dto.UserDTO{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Permissions:  permissionCodes,
	}

	return userData, nil
}

// CreateUser creates a new user in the storage and returns new user ID.
func (c *CRUD) CreateUser(ctx context.Context, user dto.CreateUserDTO) (int64, error) {
	now := time.Now()
	userToCreate := domain.User{
		Username:      user.Username,
		PasswordHash:  string(user.PasswordHash),
		FIO:           user.FIO,
		DateOfBirth:   null.TimeFromPtr(user.DateOfBirth),
		Gender:        user.Gender.ToNullInt16(),
		AvatarFileKey: null.StringFromPtr(user.AvatarFileKey),
		Deleted:       false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	newUserID, err := c.userSaver.CreateUser(ctx, userToCreate)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// CreateUserClientLink creates a user's client link and returns new link ID.
func (c *CRUD) CreateUserClientLink(ctx context.Context, userID int64, clientID int64) error {
	now := time.Now()
	userClientLinkToCreate := domain.UserClientLink{
		UserID:    userID,
		ClientID:  clientID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := c.userSaver.CreateUserClientLink(ctx, userClientLinkToCreate); err != nil {
		return err
	}

	return nil
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
