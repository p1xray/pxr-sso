package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToRoleDTO(role models.Role) dto.Role {
	permissions := make([]dto.Permission, len(role.PermissionLinks))
	for i, link := range role.PermissionLinks {
		permissions[i] = ToPermissionDTO(link.Permission)
	}

	roleDTO := dto.NewRole(role.ID, role.Code, role.Name, role.Description, permissions)
	return roleDTO
}
