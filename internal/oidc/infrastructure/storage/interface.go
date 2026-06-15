package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

type Storage interface {
	WithTransaction(ctx context.Context, f func(pgx.Tx) error) error

	Client(ctx context.Context, id int64) (models.Client, error)
	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)
	ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error)
	ClientScopeLinks(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error)
	ClientDefaultRoleLinks(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error)

	User(ctx context.Context, id int64) (models.User, error)
	IsUserExistByUsername(ctx context.Context, username string) (bool, error)
	UserByUsername(ctx context.Context, username string) (models.User, error)
	UserRoleLinks(ctx context.Context, userID int64) ([]models.UserRoleLink, error)

	CreateUser(ctx context.Context, tx pgx.Tx, user models.User) (int64, error)
	CreateUserClientLink(ctx context.Context, tx pgx.Tx, link models.UserClientLink) (int64, error)
	CreateUserRoleLinks(ctx context.Context, tx pgx.Tx, links []models.UserRoleLink) error

	RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error)

	Scopes(ctx context.Context, ids []int64) ([]models.Scope, error)
	ScopesByCode(ctx context.Context, codes []string) ([]models.Scope, error)

	Roles(ctx context.Context, ids []int64) ([]models.Role, error)
	Permissions(ctx context.Context, ids []int64) ([]models.Permission, error)

	Session(ctx context.Context, id int64) (models.Session, error)
	SessionsByCode(ctx context.Context, codes []string) ([]models.Session, error)
	SessionByCode(ctx context.Context, code string) (models.Session, error)

	SessionGrantedScopeLinks(ctx context.Context, sessionID int64) ([]models.SessionGrantedScopeLink, error)
	CreateSessionGrantedScopeLinks(ctx context.Context, tx pgx.Tx, links []models.SessionGrantedScopeLink) error
	RemoveSessionGrantedScopeLinks(ctx context.Context, tx pgx.Tx, links []int64) error
	RemoveSessionGrantedScopeLinksBySessionID(ctx context.Context, tx pgx.Tx, sessionID int64) error

	CreateSession(ctx context.Context, tx pgx.Tx, session models.Session) (int64, error)
	UpdateSession(ctx context.Context, tx pgx.Tx, session models.Session) error
	RemoveSession(ctx context.Context, tx pgx.Tx, sessionID int64) error
}
