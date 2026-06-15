package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/pkg/postgresql"
)

const pkgTag = "postgresql storage"

// storage provides access to PostgreSQL storage.
type storage struct {
	pg *postgresql.Postgres
}

// New creates a new instance of the PostgreSQL store.
func New(cfg Config) (*storage, error) {
	const op = "create new instance"

	pg, err := postgresql.New(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return &storage{pg: pg}, nil
}

// Close closes connections to PostgreSQL.
func (s *storage) Close() {
	s.pg.Close()
}

func (s *storage) WithTransaction(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, "begin transaction", err)
	}

	if err = f(tx); err != nil {
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
