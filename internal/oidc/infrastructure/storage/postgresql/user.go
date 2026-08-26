package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) User(ctx context.Context, id int64) (models.User, error) {
	const op = "get user by id"

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
		where u.id = @id;`

	args := pgx.NamedArgs{
		"id": id,
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
			return models.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return user, nil
}

func (s *storage) IsUserExistByUsername(ctx context.Context, username string) (bool, error) {
	const op = "check if user exists by username"

	stmt := `select exists(select 1 from sso.users u where u.username = @username);`

	args := pgx.NamedArgs{
		"username": username,
	}

	exists := false
	err := s.pg.Pool.QueryRow(ctx, stmt, args).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return exists, nil
}

func (s *storage) UserByUsername(ctx context.Context, username string) (models.User, error) {
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
			return models.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return user, nil
}

func (s *storage) UserRoleLinks(ctx context.Context, userID int64) ([]models.UserRoleLink, error) {
	const op = "get user role links"

	stmt :=
		`select
			 link.id,
			 link.user_id,
			 link.role_id,
			 link.created_at,
			 link.updated_at
		 from sso.user_role_links link
		 where link.user_id = @user_id;`

	args := pgx.NamedArgs{
		"user_id": userID,
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

func (s *storage) CreateUser(ctx context.Context, tx pgx.Tx, user models.User) (int64, error) {
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

	row := tx.QueryRow(ctx, stmt, args)

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

func (s *storage) CreateUserClientLink(ctx context.Context, tx pgx.Tx, link models.UserClientLink) (int64, error) {
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

	row := tx.QueryRow(ctx, stmt, args)

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

func (s *storage) CreateUserRoleLinks(ctx context.Context, tx pgx.Tx, links []models.UserRoleLink) error {
	const op = "create new user role links"

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"sso", "user_role_links"},
		[]string{"user_id", "role_id", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(links), func(i int) ([]any, error) {
			return []any{links[i].UserID, links[i].RoleID, links[i].CreatedAt, links[i].UpdatedAt}, nil
		}),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityExists, pgErr.Error())
		}

		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}
