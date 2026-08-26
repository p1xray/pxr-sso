package converter

import "github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"

func ToUserRoleLinkStorage(userID, roleID int64, setters ...models.UserRoleLinkOption) models.UserRoleLink {
	userRoleLinkModel := models.UserRoleLink{
		UserID: userID,
		RoleID: roleID,
	}

	for _, setter := range setters {
		setter(&userRoleLinkModel)
	}

	return userRoleLinkModel
}
