package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

func ToUserDTONew(user models.User) dto.User {
	roles := make([]dto.Role, len(user.RoleLinks))
	for i, link := range user.RoleLinks {
		roleDTO := ToRoleDTO(link.Role)
		roles[i] = roleDTO
	}

	userDTO := dto.NewUser(
		user.ID,
		user.Username,
		user.PasswordHash,
		roles,
	)

	return userDTO
}

func ToUserDTO(user models.User) dto.User {
	return dto.NewUser(
		user.ID,
		user.Username,
		user.PasswordHash,
		[]dto.Role{},
	)
}

func ToUserStorage(dst models.User, src dto.User, setters ...models.UserOption) models.User {
	dst.ID = src.ID()
	dst.Username = src.Username()
	dst.PasswordHash = src.PasswordHash()
	dst.FullName = src.FullName()

	for _, setter := range setters {
		setter(&dst)
	}

	return dst
}
