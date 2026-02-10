package user

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

const emptyID = 0

type Storage interface {
	WithTransaction(ctx context.Context, f func() error) error

	UserByUsername(ctx context.Context, username string) (models.User, error)
	UserRoleLinks(ctx context.Context, userID int64) ([]models.UserRoleLink, error)
	Roles(ctx context.Context, ids []int64) ([]models.Role, error)
	RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error)
	Permissions(ctx context.Context, ids []int64) ([]models.Permission, error)

	CreateUser(ctx context.Context, user models.User) (int64, error)
	CreateUserClientLink(ctx context.Context, link models.UserClientLink) (int64, error)
	CreateUserRoleLink(ctx context.Context, link models.UserRoleLink) (int64, error)
}

type Repository struct {
	storage Storage
}

func NewRepository(storage Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

type Option func(context.Context, *Repository, *models.User) error

func (r *Repository) UserByUsername(ctx context.Context, username string, opts ...Option) (dto.User, error) {
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

	userDTO := converter.ToUserDTONew(user)
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

func (r *Repository) roles(ctx context.Context, ids []int64) ([]models.Role, error) {
	roles, err := r.storage.Roles(ctx, ids)
	if err != nil {
		return []models.Role{}, fmt.Errorf("%s: %w", "get roles", err)
	}

	rolePermissionLinks, err := r.rolePermissions(ctx, ids)
	if err != nil {
		return []models.Role{}, err
	}

	for i := range roles {
		rolePermissions := make([]models.RolePermissionLink, 0)
		for j := range rolePermissionLinks {
			if roles[i].ID == rolePermissionLinks[j].RoleID {
				rolePermissions = append(rolePermissions, rolePermissionLinks[j])
			}
		}

		roles[i].PermissionLinks = rolePermissions
	}

	return roles, nil
}

func (r *Repository) rolePermissions(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error) {
	rolePermissionLinks, err := r.storage.RolePermissionLinks(ctx, roleIDs)
	if err != nil {
		return []models.RolePermissionLink{}, fmt.Errorf("%s: %w", "get role permission links", err)
	}

	permissionIDs := make([]int64, len(rolePermissionLinks))
	for i, link := range rolePermissionLinks {
		permissionIDs[i] = link.PermissionID
	}

	permissions, err := r.storage.Permissions(ctx, permissionIDs)
	if err != nil {
		return []models.RolePermissionLink{}, fmt.Errorf("%s: %w", "get permissions", err)
	}

	for i := range rolePermissionLinks {
		for j := range permissions {
			if rolePermissionLinks[i].PermissionID == permissions[j].ID {
				rolePermissionLinks[i].Permission = permissions[j]
				break
			}
		}
	}

	return rolePermissionLinks, nil
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

func WithRoles() Option {
	return func(ctx context.Context, r *Repository, user *models.User) error {
		userRoleLinks, err := r.userRoles(ctx, user.ID)
		if err != nil {
			return err
		}

		user.RoleLinks = userRoleLinks
		return nil
	}
}
