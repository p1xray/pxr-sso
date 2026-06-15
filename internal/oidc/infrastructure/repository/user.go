package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/enum"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

type UserOption func(context.Context, *repository, *models.User) error

func (r *repository) User(ctx context.Context, id int64, opts ...UserOption) (dto.User, error) {
	const op = "get user by id"

	user, err := r.storage.User(ctx, id)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &user); err != nil {
			return dto.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
	}

	userDTO := converter.ToUserDTO(user)
	return userDTO, nil
}

func (r *repository) UserByUsername(ctx context.Context, username string, opts ...UserOption) (dto.User, error) {
	const op = "get user by username"

	user, err := r.storage.UserByUsername(ctx, username)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &user); err != nil {
			return dto.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
	}

	userDTO := converter.ToUserDTO(user)
	return userDTO, nil
}

func (r *repository) IsUserExistByUsername(ctx context.Context, username string) (bool, error) {
	const op = "check user exists by username"

	isUserExist, err := r.storage.IsUserExistByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return isUserExist, nil
}

func (r *repository) userRoles(ctx context.Context, userID int64) ([]models.UserRoleLink, error) {
	userRoleLinks, err := r.storage.UserRoleLinks(ctx, userID)
	if err != nil {
		return []models.UserRoleLink{}, err
	}

	roleIDs := make([]int64, len(userRoleLinks))
	for i, link := range userRoleLinks {
		roleIDs[i] = link.RoleID
	}

	roles, err := r.roles(ctx, roleIDs)
	if err != nil {
		return []models.UserRoleLink{}, err
	}

	for i := range userRoleLinks {
		for j := range roles {
			if userRoleLinks[i].RoleID == roles[j].ID {
				userRoleLinks[i].Role = roles[j]
				break
			}
		}
	}

	return userRoleLinks, nil
}

func (r *repository) SaveUser(ctx context.Context, user entity.User) (dto.User, error) {
	const op = "save user"

	savedUser, err := r.saveUserByDataStatus(ctx, user)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return savedUser, nil
}

func (r *repository) saveUserByDataStatus(ctx context.Context, user entity.User) (dto.User, error) {
	switch user.DataStatus() {
	case enum.DataStatusToCreate:
		userID, err := r.createUser(ctx, user)
		if err != nil {
			return dto.User{}, err
		}

		savedUser, err := r.User(ctx, userID)
		if err != nil {
			return dto.User{}, err
		}

		return savedUser, nil
	default:
		return dto.User{}, fmt.Errorf(
			"there is no implementation of save user for data status with value: %d", user.DataStatus())
	}
}

func (r *repository) createUser(ctx context.Context, user entity.User) (int64, error) {
	userToCreate := converter.ToUserStorage(user, models.UserCreated())

	newUserID := int64(0)
	err := r.storage.WithTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		newUserID, err = r.storage.CreateUser(ctx, tx, userToCreate)
		if err != nil {
			return err
		}

		userClientLinkToCreate := converter.ToUserClientLinkStorage(
			newUserID,
			user.ClientID(),
			models.UserClientLinkCreated(),
		)
		if _, err = r.storage.CreateUserClientLink(ctx, tx, userClientLinkToCreate); err != nil {
			return err
		}

		userRoleLinksToCreate := make([]models.UserRoleLink, 0)
		for _, role := range user.Roles() {
			userRoleLinkToCreate := converter.ToUserRoleLinkStorage(newUserID, role.ID(), models.UserRoleLinkCreated())
			userRoleLinksToCreate = append(userRoleLinksToCreate, userRoleLinkToCreate)
		}

		if err = r.storage.CreateUserRoleLinks(ctx, tx, userRoleLinksToCreate); err != nil {
			return err
		}

		return nil
	})

	return newUserID, err
}

func WithRoles() UserOption {
	return func(ctx context.Context, r *repository, user *models.User) error {
		userRoleLinks, err := r.userRoles(ctx, user.ID)
		if err != nil {
			return err
		}

		user.RoleLinks = userRoleLinks
		return nil
	}
}
