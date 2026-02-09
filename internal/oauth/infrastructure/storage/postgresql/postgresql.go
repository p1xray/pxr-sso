package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/pkg/postgresql"
)

// Storage provides access to PostgreSQL storage.
type Storage struct {
	pg *postgresql.Postgres
}

// New creates a new instance of the PostgreSQL store.
func New(connectionURL string) (*Storage, error) {
	const op = "infrastructure.storage.postgresql.New"

	pg, err := postgresql.New(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return err
	}

	if err = f(); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Storage) ClientByCode(ctx context.Context, code string) (models.Client, error) {
	const op = "infrastructure.storage.postgresql.ClientByCode"

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
			return models.Client{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

func (s *Storage) ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error) {
	const op = "infrastructure.storage.postgresql.ClientAudiences"

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
		return nil, fmt.Errorf("%s: %w", op, err)
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		audiences = append(audiences, audience)
	}

	return audiences, nil
}

func (s *Storage) ClientRedirectURIs(ctx context.Context, clientID int64) ([]models.RedirectURI, error) {
	const op = "infrastructure.storage.postgresql.ClientRedirectURIs"

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
		return nil, fmt.Errorf("%s: %w", op, err)
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		redirectURIs = append(redirectURIs, redirectURI)
	}

	return redirectURIs, nil
}

func (s *Storage) ClientScopes(ctx context.Context, clientID int64) ([]models.Scope, error) {
	const op = "infrastructure.storage.postgresql.ClientScopes"

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
		return nil, fmt.Errorf("%s: %w", op, err)
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "infrastructure.storage.postgresql.UserByUsername"

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
			return models.User{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const op = "infrastructure.storage.postgresql.CreateUser"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateUserClientLink(ctx context.Context, link models.UserClientLink) (int64, error) {
	const op = "infrastructure.storage.postgresql.CreateUserClientLink"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateUserRoleLink(ctx context.Context, link models.UserRoleLink) (int64, error) {
	const op = "infrastructure.storage.postgresql.CreateUserRoleLink"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
