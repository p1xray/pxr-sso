package converter

import (
	"github.com/guregu/null/v6"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToUserDTO(user models.User) dto.User {
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

func ToUserStorage(user entity.User, setters ...models.UserOption) models.User {
	dateOfBirth := user.DateOfBirth()
	gender := user.Gender()
	avatarFileKey := user.AvatarFileKey()

	userStorageModel := models.User{
		Username:      user.Username(),
		PasswordHash:  user.PasswordHash(),
		FullName:      user.FullName(),
		DateOfBirth:   null.NewTime(dateOfBirth.Unwrap(), dateOfBirth.IsSome()),
		Gender:        null.NewInt16(int16(gender.Unwrap()), gender.IsSome()),
		AvatarFileKey: null.NewString(avatarFileKey.Unwrap(), avatarFileKey.IsSome()),
	}

	for _, setter := range setters {
		setter(&userStorageModel)
	}

	return userStorageModel
}
