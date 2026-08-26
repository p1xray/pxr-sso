package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) Client(ctx context.Context, id int64) (models.Client, error) {
	const op = "get client by id"

	stmt :=
		`select
			 c.id,
			 c.name,
			 c.code,
			 c.secret_key,
			 c.deleted,
			 c.created_at,
			 c.updated_at
		 from sso.clients c
		 where c.id = @id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var client models.Client
	err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Code,
		&client.SecretKey,
		&client.Deleted,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Client{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.Client{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return client, nil
}

func (s *storage) ClientByCode(ctx context.Context, code string) (models.Client, error) {
	const op = "get client by code"

	stmt :=
		`select
			 c.id,
			 c.name,
			 c.code,
			 c.secret_key,
			 c.deleted,
			 c.created_at,
			 c.updated_at
		 from sso.clients c
		 where c.code = @code;`

	args := pgx.NamedArgs{
		"code": code,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var client models.Client
	err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Code,
		&client.SecretKey,
		&client.Deleted,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Client{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.Client{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return client, nil
}

func (s *storage) ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error) {
	const op = "get client audiences"

	stmt :=
		`select
			 ca.id,
			 ca.client_id,
			 ca.uri,
			 ca.created_at,
			 ca.updated_at
		 from sso.client_audiences ca
		 where ca.client_id = @client_id;`

	args := pgx.NamedArgs{
		"client_id": clientID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	audiences := make([]models.Audience, 0)
	for rows.Next() {
		audience := models.Audience{}
		err = rows.Scan(
			&audience.ID,
			&audience.ClientID,
			&audience.URI,
			&audience.CreatedAt,
			&audience.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		audiences = append(audiences, audience)
	}

	return audiences, nil
}

func (s *storage) ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error) {
	const op = "get client redirect uris"

	stmt :=
		`select
			 cri.id,
			 cri.client_id,
			 cri.uri,
			 cri.created_at,
			 cri.updated_at
		 from sso.client_redirect_uris cri
		 where cri.client_id = @client_id;`

	args := pgx.NamedArgs{
		"client_id": clientID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	redirectURIs := make([]models.RedirectURI, 0)
	for rows.Next() {
		redirectURI := models.RedirectURI{}
		err = rows.Scan(
			&redirectURI.ID,
			&redirectURI.ClientID,
			&redirectURI.URI,
			&redirectURI.CreatedAt,
			&redirectURI.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		redirectURIs = append(redirectURIs, redirectURI)
	}

	return redirectURIs, nil
}

func (s *storage) ClientScopeLinks(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error) {
	const op = "get client scope links"

	stmt :=
		`select
			 link.id,
			 link.client_id,
			 link.scope_id,
			 link.created_at,
			 link.updated_at
		 from sso.client_scope_links link
		 where link.client_id = @client_id;`

	args := pgx.NamedArgs{
		"client_id": clientID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	links := make([]models.ClientScopeLink, 0)
	for rows.Next() {
		link := models.ClientScopeLink{}
		err = rows.Scan(
			&link.ID,
			&link.ClientID,
			&link.ScopeID,
			&link.CreatedAt,
			&link.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		links = append(links, link)
	}

	return links, nil
}

func (s *storage) ClientDefaultRoleLinks(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error) {
	const op = "get client default role links"

	stmt :=
		`select
			 link.id,
			 link.client_id,
			 link.role_id,
			 link.created_at,
			 link.updated_at
		 from sso.client_default_role_links link
		 where link.client_id = @client_id;`

	args := pgx.NamedArgs{
		"client_id": clientID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	links := make([]models.ClientDefaultRoleLink, 0)
	for rows.Next() {
		link := models.ClientDefaultRoleLink{}
		err = rows.Scan(
			&link.ID,
			&link.ClientID,
			&link.RoleID,
			&link.CreatedAt,
			&link.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		links = append(links, link)
	}

	return links, nil
}

func (s *storage) ClientScopes(ctx context.Context, clientID int64) ([]models.Scope, error) {
	const op = "get client scopes"

	stmt :=
		`select
			 s.id,
			 s.code,
			 s.name,
			 s.description,
			 s.created_at,
			 s.updated_at
		 from sso.scopes s
		 	join sso.client_scope_links csl on csl.scope_id = s.id
		 where csl.client_id = @client_id;`

	args := pgx.NamedArgs{
		"client_id": clientID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	scopes := make([]models.Scope, 0)
	for rows.Next() {
		scope := models.Scope{}
		err = rows.Scan(
			&scope.ID,
			&scope.Code,
			&scope.Name,
			&scope.Description,
			&scope.CreatedAt,
			&scope.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}
