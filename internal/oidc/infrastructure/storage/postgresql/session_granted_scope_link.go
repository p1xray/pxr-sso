package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) SessionGrantedScopeLinks(ctx context.Context, sessionID int64) ([]models.SessionGrantedScopeLink, error) {
	const op = "get session granted scope links by session id"

	stmt :=
		`select
			 link.id,
			 link.session_id,
			 link.scope_id,
			 link.created_at,
			 link.updated_at
		 from sso.session_granted_scope_links link
		 where link.session_id = @session_id;`

	args := pgx.NamedArgs{
		"session_id": sessionID,
	}

	rows, err := s.pg.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}
	defer rows.Close()

	links := make([]models.SessionGrantedScopeLink, 0)
	for rows.Next() {
		link := models.SessionGrantedScopeLink{}
		err = rows.Scan(
			&link.ID,
			&link.SessionID,
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

func (s *storage) CreateSessionGrantedScopeLinks(ctx context.Context, links []models.SessionGrantedScopeLink) error {
	const op = "create new session granted scope links"

	_, err := s.pg.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"sso", "session_granted_scope_links"},
		[]string{"session_id", "scope_id", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(links), func(i int) ([]any, error) {
			return []any{links[i].SessionID, links[i].ScopeID, links[i].CreatedAt, links[i].UpdatedAt}, nil
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

func (s *storage) RemoveSessionGrantedScopeLinks(ctx context.Context, links []int64) error {
	const op = "remove session granted scope links"

	stmt := `delete from sso.session_granted_scope_links where id = any(@ids);`

	args := pgx.NamedArgs{
		"ids": links,
	}

	_, err := s.pg.Pool.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (s *storage) RemoveSessionGrantedScopeLinksBySessionID(ctx context.Context, sessionID int64) error {
	const op = "remove session granted scope links by session id"

	stmt := `delete from sso.session_granted_scope_links where session_id = @session_id;`

	args := pgx.NamedArgs{
		"session_id": sessionID,
	}

	_, err := s.pg.Pool.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}
