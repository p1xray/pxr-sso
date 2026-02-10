package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

func ToRoleDTO(role models.Role) dto.Role {
	permissions := make([]string, len(role.PermissionLinks))
	for i, link := range role.PermissionLinks {
		permissions[i] = link.Permission.Code
	}

	roleDTO := dto.NewRole(role.ID, role.Code, permissions)
	return roleDTO
}
