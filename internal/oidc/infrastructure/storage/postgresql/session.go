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

func (s *storage) Session(ctx context.Context, id int64) (models.Session, error) {
	const op = "get session by id"

	stmt :=
		`select
			 s.id,
			 s.code,
			 s.client_id,
			 s.user_id,
			 s.auth_time,
			 s.identity_provider,
			 s.created_at,
			 s.updated_at
		 from sso.sessions s
		 where s.id = @id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.Code,
		&session.ClientID,
		&session.UserID,
		&session.AuthTime,
		&session.IdentityProvider,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return session, nil
}

func (s *storage) SessionsByCode(ctx context.Context, codes []string) ([]models.Session, error) {
	const op = "get session by code"

	stmt :=
		`select
			 s.id,
			 s.code,
			 s.client_id,
			 s.user_id,
			 s.auth_time,
			 s.identity_provider,
			 s.created_at,
			 s.updated_at
		 from sso.sessions s
		 where s.code = any(@codes);`

	args := pgx.NamedArgs{
		"codes": codes,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	sessions := make([]models.Session, 0)
	for rows.Next() {
		session := models.Session{}
		err = rows.Scan(
			&session.ID,
			&session.Code,
			&session.ClientID,
			&session.UserID,
			&session.AuthTime,
			&session.IdentityProvider,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *storage) SessionByCode(ctx context.Context, code string) (models.Session, error) {
	const op = "get session by code"

	stmt :=
		`select
			 s.id,
			 s.code,
			 s.client_id,
			 s.user_id,
			 s.auth_time,
			 s.identity_provider,
			 s.created_at,
			 s.updated_at
		 from sso.sessions s
		 where s.code = @code;`

	args := pgx.NamedArgs{
		"code": code,
	}

	row := s.pg.Pool.QueryRow(ctx, stmt, args)

	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.Code,
		&session.ClientID,
		&session.UserID,
		&session.AuthTime,
		&session.IdentityProvider,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return models.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return session, nil
}

func (s *storage) CreateSession(ctx context.Context, tx pgx.Tx, session models.Session) (int64, error) {
	const op = "create new session"

	stmt :=
		`insert into sso.sessions (
		   code,
		   client_id,
		   user_id,
		   auth_time,
		   identity_provider,
		   created_at,
		   updated_at)
		values(
		   @code,
		   @client_id,
		   @user_id,
		   @auth_time,
		   @identity_provider,
		   @created_at,
		   @updated_at)
		returning id;`

	args := pgx.NamedArgs{
		"code":              session.Code.String(),
		"client_id":         session.ClientID,
		"user_id":           session.UserID,
		"auth_time":         session.AuthTime,
		"identity_provider": session.IdentityProvider,
		"created_at":        session.CreatedAt,
		"updated_at":        session.UpdatedAt,
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

func (s *storage) UpdateSession(ctx context.Context, tx pgx.Tx, session models.Session) error {
	const op = "update session"

	stmt :=
		`update sso.sessions
		set code = @code,
		    client_id = @client_id,
		    user_id = @user_id,
		    auth_time = @auth_time,
		    identity_provider = @identity_provider,
		    updated_at = @updated_at
		where id = @id;`

	args := pgx.NamedArgs{
		"id":                session.ID,
		"code":              session.Code,
		"client_id":         session.ClientID,
		"user_id":           session.UserID,
		"auth_time":         session.AuthTime,
		"identity_provider": session.IdentityProvider,
		"updated_at":        session.UpdatedAt,
	}

	_, err := tx.Exec(ctx, stmt, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s: %s: %w: %s", pkgTag, op, infrastructure.ErrEntityExists, pgErr.Error())
		}

		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (s *storage) RemoveSession(ctx context.Context, tx pgx.Tx, sessionID int64) error {
	const op = "remove session"

	stmt := `delete from sso.sessions where id = @id;`

	args := pgx.NamedArgs{
		"id": sessionID,
	}

	_, err := tx.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}
