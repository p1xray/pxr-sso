package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) Permissions(ctx context.Context, ids []int64) ([]models.Permission, error) {
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
