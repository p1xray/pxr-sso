package repository

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

type ClientOption func(context.Context, *Repository, *models.Client) error

func (r *Repository) ClientByCode(ctx context.Context, code string, opts ...ClientOption) (dto.Client, error) {
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

	clientDTO := converter.ToClientDTO(client)
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

func WithAudiences() ClientOption {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		audiences, err := r.clientAudiences(ctx, client.ID)
		if err != nil {
			return err
		}

		client.Audiences = audiences
		return nil
	}
}

func WithRedirectURIs() ClientOption {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		redirectURIs, err := r.clientRedirectURIs(ctx, client.ID)
		if err != nil {
			return err
		}

		client.RedirectURIs = redirectURIs
		return nil
	}
}

func WithScopes() ClientOption {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		clientScopeLinks, err := r.clientScopes(ctx, client.ID)
		if err != nil {
			return err
		}

		client.ScopeLinks = clientScopeLinks
		return nil
	}
}

func WithDefaultRoles() ClientOption {
	return func(ctx context.Context, r *Repository, client *models.Client) error {
		clientDefaultRoleLinks, err := r.clientDefaultRoles(ctx, client.ID)
		if err != nil {
			return err
		}

		client.DefaultRoleLinks = clientDefaultRoleLinks
		return nil
	}
}
