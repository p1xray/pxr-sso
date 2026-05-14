package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) Scopes(ctx context.Context, ids []int64) ([]models.Scope, error) {
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
