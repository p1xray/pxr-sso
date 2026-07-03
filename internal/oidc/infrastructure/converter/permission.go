package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToPermissionDTO(permission models.Permission) dto.Permission {
	return dto.NewPermission(permission.ID, permission.Code, permission.Description)
}
