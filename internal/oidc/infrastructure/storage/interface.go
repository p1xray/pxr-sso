package storage

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

type Storage interface {
	WithTransaction(ctx context.Context, f func() error) error

	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)
	ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error)
	ClientScopeLinks(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error)
	ClientDefaultRoleLinks(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error)

	UserByUsername(ctx context.Context, username string) (models.User, error)
	UserRoleLinks(ctx context.Context, userID int64) ([]models.UserRoleLink, error)

	CreateUser(ctx context.Context, user models.User) (int64, error)
	CreateUserClientLink(ctx context.Context, link models.UserClientLink) (int64, error)
	CreateUserRoleLink(ctx context.Context, link models.UserRoleLink) (int64, error)

	RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error)

	Scopes(ctx context.Context, ids []int64) ([]models.Scope, error)
	Roles(ctx context.Context, ids []int64) ([]models.Role, error)
	Permissions(ctx context.Context, ids []int64) ([]models.Permission, error)
}
