package client

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

type Storage interface {
	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)
	ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error)
	ClientScopeLinks(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error)
	Scopes(ctx context.Context, ids []int64) ([]models.Scope, error)
	ClientDefaultRoleLinks(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error)
	Roles(ctx context.Context, ids []int64) ([]models.Role, error)
	RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error)
	Permissions(ctx context.Context, ids []int64) ([]models.Permission, error)
}

type Repository struct {
	storage Storage
}

func NewRepository(storage Storage) *Repository {
	return &Repository{
		storage: storage,
	}
}

type Option func(context.Context, *Repository, *models.Client) error

func (r *Repository) ClientByCode(ctx context.Context, code string, opts ...Option) (dto.Client, error) {
	const op = "client repository: get client by code"

	client, err := r.storage.ClientByCode(ctx, code)
	if err != nil {
		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &client); err != nil {
			return dto.Client{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	clientDTO := converter.ToClientDTONew(client)
	return clientDTO, nil
}

func (r *Repository) clientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error) {
	audiences, err := r.storage.ClientAudiences(ctx, clientID)
	if err != nil {
		return []models.Audience{}, fmt.Errorf("%s: %w", "get client audiences", err)
	}

	return audiences, nil
}

func (r *Repository) clientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error) {
	redirectURIs, err := r.storage.ClientRedirectURIs(ctx, clientID)
	if err != nil {
		return []models.RedirectURI{}, fmt.Errorf("%s: %w", "get client redirect URIs", err)
	}

	return redirectURIs, nil
}

func (r *Repository) clientScopes(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error) {
	clientScopeLinks, err := r.storage.ClientScopeLinks(ctx, clientID)
	if err != nil {
		return []models.ClientScopeLink{}, fmt.Errorf("%s: %w", "get client scope links", err)
	}

	scopeIDs := make([]int64, len(clientScopeLinks))
	for i, link := range clientScopeLinks {
		scopeIDs[i] = link.ScopeID
	}

	scopes, err := r.storage.Scopes(ctx, scopeIDs)
	if err != nil {
		return []models.ClientScopeLink{}, fmt.Errorf("%s: %w", "get scopes", err)
	}

	for i := range clientScopeLinks {
		for j := range scopes {
			if clientScopeLinks[i].ScopeID == scopes[j].ID {
				clientScopeLinks[i].Scope = scopes[j]
				break
			}
		}
	}

	return clientScopeLinks, nil
}

func (r *Repository) clientDefaultRoles(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error) {
	clientDefaultRoleLinks, err := r.storage.ClientDefaultRoleLinks(ctx, clientID)
	if err != nil {
		return []models.ClientDefaultRoleLink{}, fmt.Errorf("%s: %w", "get client default role links", err)
	}

	roleIDs := make([]int64, len(clientDefaultRoleLinks))
	for i, link := range clientDefaultRoleLinks {
		roleIDs[i] = link.RoleID
	}

	roles, err := r.roles(ctx, roleIDs)
	if err != nil {
		return []models.ClientDefaultRoleLink{}, err
	}

	for i := range clientDefaultRoleLinks {
		for j := range roles {
			if clientDefaultRoleLinks[i].RoleID == roles[j].ID {
				clientDefaultRoleLinks[i].Role = roles[j]
				break
			}
		}
	}

	return clientDefaultRoleLinks, nil
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

func WithAudiences() Option {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		audiences, err := r.clientAudiences(ctx, client.ID)
		if err != nil {
			return err
		}

		client.Audiences = audiences
		return nil
	}
}

func WithRedirectURIs() Option {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		redirectURIs, err := r.clientRedirectURIs(ctx, client.ID)
		if err != nil {
			return err
		}

		client.RedirectURIs = redirectURIs
		return nil
	}
}

func WithScopes() Option {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		clientScopeLinks, err := r.clientScopes(ctx, client.ID)
		if err != nil {
			return err
		}

		client.ScopeLinks = clientScopeLinks
		return nil
	}
}

func WithDefaultRoles() Option {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		clientDefaultRoleLinks, err := r.clientDefaultRoles(ctx, client.ID)
		if err != nil {
			return err
		}

		client.DefaultRoleLinks = clientDefaultRoleLinks
		return nil
	}
}
