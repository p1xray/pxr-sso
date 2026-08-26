package repository

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func (r *repository) roles(ctx context.Context, ids []int64) ([]models.Role, error) {
	roles, err := r.storage.Roles(ctx, ids)
	if err != nil {
		return []models.Role{}, err
	}

	rolePermissionLinks, err := r.rolePermissions(ctx, ids)
	if err != nil {
		return []models.Role{}, err
	}

	for i := range roles {
		rolePermissions := make([]models.RolePermissionLink, 0)
		for j := range rolePermissionLinks {
			if roles[i].ID == rolePermissionLinks[j].RoleID {
				rolePermissions = append(rolePermissions, rolePermissionLinks[j])
			}
		}

		roles[i].PermissionLinks = rolePermissions
	}

	return roles, nil
}

func (r *repository) rolePermissions(ctx context.Context, roleIDs []int64) ([]models.RolePermissionLink, error) {
	rolePermissionLinks, err := r.storage.RolePermissionLinks(ctx, roleIDs)
	if err != nil {
		return []models.RolePermissionLink{}, err
	}

	permissionIDs := make([]int64, len(rolePermissionLinks))
	for i, link := range rolePermissionLinks {
		permissionIDs[i] = link.PermissionID
	}

	permissions, err := r.storage.Permissions(ctx, permissionIDs)
	if err != nil {
		return []models.RolePermissionLink{}, err
	}

	for i := range rolePermissionLinks {
		for j := range permissions {
			if rolePermissionLinks[i].PermissionID == permissions[j].ID {
				rolePermissionLinks[i].Permission = permissions[j]
				break
			}
		}
	}

	return rolePermissionLinks, nil
}
