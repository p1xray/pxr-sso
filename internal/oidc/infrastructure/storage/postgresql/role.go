package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (s *storage) Roles(ctx context.Context, ids []int64) ([]models.Role, error) {
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

func (s *storage) RolePermissionLinks(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error) {
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
