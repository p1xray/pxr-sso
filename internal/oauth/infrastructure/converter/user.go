package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

func ToUserDTO(user models.User) dto.User {
	return dto.NewUser(
		user.ID,
		user.Username,
		user.PasswordHash,
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
