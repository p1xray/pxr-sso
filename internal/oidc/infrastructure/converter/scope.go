package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToScopeDTO(scope models.Scope) dto.Scope {
	return dto.NewScope(scope.ID, scope.Code, scope.Name, scope.Description)
}
