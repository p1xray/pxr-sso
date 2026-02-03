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
