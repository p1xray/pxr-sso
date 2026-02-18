package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/pkg/postgresql"
)

const pkgTag = "postgresql storage"

// Storage provides access to PostgreSQL storage.
type Storage struct {
	pg *postgresql.Postgres
}

// New creates a new instance of the PostgreSQL store.
func New(cfg Config) (*Storage, error) {
	const op = "create new instance"

	pg, err := postgresql.New(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return &Storage{pg: pg}, nil
}

// Close closes connections to PostgreSQL.
func (s *Storage) Close() {
	s.pg.Close()
}

func (s *Storage) WithTransaction(ctx context.Context, f func() error) error {
	tx, err := s.pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, "begin transaction", err)
	}

	if err = f(); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("%s: %s: %w", pkgTag, "rollback transaction", rbErr)
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, "commit transaction", err)
	}

	return nil
}

func (s *Storage) ClientByCode(ctx context.Context, code string) (models.Client, error) {
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
			return models.Client{}, fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityNotFound, err.Error())
		}

		return models.Client{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return client, nil
}

func (s *Storage) ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error) {
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

func (s *Storage) ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error) {
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

func (s *Storage) ClientScopeLinks(ctx context.Context, clientID int64) ([]models.ClientScopeLink, error) {
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

func (s *Storage) Scopes(ctx context.Context, ids []int64) ([]models.Scope, error) {
	const op = "get scopes"

	stmt :=
		`select
			 s.id,
			 s.code,
			 s.name,
			 s.description,
			 s.created_at,
			 s.updated_at
		 from sso.scopes s
		 where s.id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": ids,
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

func (s *Storage) ClientDefaultRoleLinks(ctx context.Context, clientID int64) ([]models.ClientDefaultRoleLink, error) {
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

func (s *Storage) Roles(ctx context.Context, ids []int64) ([]models.Role, error) {
	const op = "get roles"

	stmt :=
		`select
			 r.id,
			 r.code,
			 r.name,
			 r.description,
			 r.active,
			 r.deleted,
			 r.created_at,
			 r.updated_at
		 from sso.roles r
		 where r.id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": ids,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		role := models.Role{}
		err = rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.Active,
			&role.Deleted,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Storage) RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error) {
	const op = "get role permission links"

	stmt :=
		`select
			 link.id,
			 link.role_id,
			 link.permission_id,
			 link.created_at,
			 link.updated_at
		 from sso.role_permission_links link
		 where link.role_id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": roleIDs,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	links := make([]models.RolePermissionLink, 0)
	for rows.Next() {
		link := models.RolePermissionLink{}
		err = rows.Scan(
			&link.ID,
			&link.RoleID,
			&link.PermissionID,
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

func (s *Storage) Permissions(ctx context.Context, ids []int64) ([]models.Permission, error) {
	const op = "get permissions"

	stmt :=
		`select
			 p.id,
			 p.code,
			 p.description,
			 p.active,
			 p.deleted,
			 p.created_at,
			 p.updated_at
		 from sso.permissions p
		 where p.id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": ids,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	permissions := make([]models.Permission, 0)
	for rows.Next() {
		permission := models.Permission{}
		err = rows.Scan(
			&permission.ID,
			&permission.Code,
			&permission.Description,
			&permission.Active,
			&permission.Deleted,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (s *Storage) ClientScopes(ctx context.Context, clientID int64) ([]models.Scope, error) {
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

func (s *Storage) UserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "get user by username"

	stmt :=
		`select
    		u.id,
    		u.username,
    		u.password_hash,
    		u.full_name,
    		u.date_of_birth,
    		u.gender,
    		u.avatar_file_key,
    		u.deleted,
    		u.created_at,
    		u.updated_at
		from sso.users u
		where u.username = @username;`

	args := pgx.NamedArgs{
		"username": username,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityNotFound, err.Error())
		}

		return models.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return user, nil
}

func (s *Storage) UserRoleLinks(ctx context.Context, userID int64) ([]models.UserRoleLink, error) {
	const op = "get user role links"

	stmt :=
		`select
			 link.id,
			 link.user_id,
			 link.role_id,
			 link.created_at,
			 link.updated_at
		 from sso.user_role_links link
		 where link.role_id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": userID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	links := make([]models.UserRoleLink, 0)
	for rows.Next() {
		link := models.UserRoleLink{}
		err = rows.Scan(
			&link.ID,
			&link.UserID,
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

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const op = "create new user"

	stmt :=
		`insert into sso.users (
		   username,
		   password_hash,
		   full_name,
		   date_of_birth,
		   gender,
		   avatar_file_key,
		   deleted,
		   created_at,
		   updated_at)
		values(
		   @username,
		   @password_hash,
		   @full_name,
		   @date_of_birth,
		   @gender,
		   @avatar_file_key,
		   @deleted,
		   @created_at,
		   @updated_at)
		returning id;`

	args := pgx.NamedArgs{
		"username":        user.Username,
		"password_hash":   user.PasswordHash,
		"full_name":       user.FullName,
		"date_of_birth":   user.DateOfBirth,
		"gender":          user.Gender,
		"avatar_file_key": user.AvatarFileKey,
		"deleted":         user.Deleted,
		"created_at":      user.CreatedAt,
		"updated_at":      user.UpdatedAt,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityExists, pgErr.Error())
		}

		return 0, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return id, nil
}

func (s *Storage) CreateUserClientLink(ctx context.Context, link models.UserClientLink) (int64, error) {
	const op = "create new user client link"

	stmt :=
		`insert into sso.user_client_links (user_id, client_id, created_at, updated_at)
		 values (@user_id, @client_id, @created_at, @updated_at)
		 returning id;`

	args := pgx.NamedArgs{
		"user_id":    link.UserID,
		"client_id":  link.ClientID,
		"created_at": link.CreatedAt,
		"updated_at": link.UpdatedAt,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityExists, pgErr.Error())
		}

		return 0, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return id, nil
}

func (s *Storage) CreateUserRoleLink(ctx context.Context, link models.UserRoleLink) (int64, error) {
	const op = "create new user role link"

	stmt :=
		`insert into sso.user_role_links (user_id, role_id, created_at, updated_at)
		 values (@user_id, @role_id, @created_at, @updated_at)
		 returning id;`

	args := pgx.NamedArgs{
		"user_id":    link.UserID,
		"role_id":    link.RoleID,
		"created_at": link.CreatedAt,
		"updated_at": link.UpdatedAt,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityExists, pgErr.Error())
		}

		return 0, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return id, nil
}
