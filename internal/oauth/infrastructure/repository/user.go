package repository

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

type UserOption func(context.Context, *Repository, *models.User) error

func (r *Repository) UserByUsername(ctx context.Context, username string, opts ...UserOption) (dto.User, error) {
	const op = "user repository: get user by username"

	user, err := r.storage.UserByUsername(ctx, username)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &user); err != nil {
			return dto.User{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	userDTO := converter.ToUserDTO(user)
	return userDTO, nil
}

func (r *Repository) userRoles(ctx context.Context, userID int64) ([]models.UserRoleLink, error) {
	userRoleLinks, err := r.storage.UserRoleLinks(ctx, userID)
	if err != nil {
		return []models.UserRoleLink{}, fmt.Errorf("%s: %w", "get user role links", err)
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

func (r *Repository) CreateUser(ctx context.Context, user dto.User, clientID int64) error {
	const op = "user repository: create user"

	err := r.storage.WithTransaction(ctx, func() error {
		newUserID, err := r.createUser(ctx, user)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = r.createUserClientLink(ctx, newUserID, clientID); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		for _, role := range user.Roles() {
			if err = r.createUserRoleLink(ctx, newUserID, role.ID()); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) createUser(ctx context.Context, user dto.User) (int64, error) {
	userStorageModel := converter.ToUserStorage(models.User{}, user, models.UserCreated())
	id, err := r.storage.CreateUser(ctx, userStorageModel)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) createUserClientLink(ctx context.Context, userID int64, clientID int64) error {
	if userID == emptyID || clientID == emptyID {
		return fmt.Errorf("a non-null identifiers is required to create an user client link in storage")
	}

	userClientLinkStorageModel := converter.ToUserClientLinkStorage(userID, clientID, models.UserClientLinkCreated())
	_, err := r.storage.CreateUserClientLink(ctx, userClientLinkStorageModel)

	return err
}

func (r *Repository) createUserRoleLink(ctx context.Context, userID, roleID int64) error {
	if userID == emptyID || roleID == emptyID {
		return fmt.Errorf("a non-null identifiers is required to create an user role link in storage")
	}

	userRoleLinkStorageModel := converter.ToUserRoleLinkStorage(userID, roleID, models.UserRoleLinkCreated())
	_, err := r.storage.CreateUserRoleLink(ctx, userRoleLinkStorageModel)

	return err
}

func WithRoles() UserOption {
	return func(ctx context.Context, r *Repository, user *models.User) error {
		userRoleLinks, err := r.userRoles(ctx, user.ID)
		if err != nil {
			return err
		}

		user.RoleLinks = userRoleLinks
		return nil
	}
}
