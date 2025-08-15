package converter

import (
	"github.com/guregu/null/v6"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/enum"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
)

func ToUserDTO(user models.User, roles []models.Role, permissions []models.Permission) dto.User {
	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	permissionCodes := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionCodes[i] = permission.Code
	}

	return dto.User{
		ID:            user.ID,
		Username:      user.Username,
		PasswordHash:  user.PasswordHash,
		FullName:      user.FullName,
		DateOfBirth:   user.DateOfBirth.Ptr(),
		Gender:        enum.GenderEnumFromNullInt16(user.Gender),
		AvatarFileKey: user.AvatarFileKey.Ptr(),
		Roles:         roleCodes,
		Permissions:   permissionCodes,
	}
}

func ToUserProfileDTO(user models.User) dto.UserProfile {
	return dto.UserProfile{
		ID:            user.ID,
		Username:      user.Username,
		FullName:      user.FullName,
		DateOfBirth:   user.DateOfBirth.Ptr(),
		Gender:        enum.GenderEnumFromNullInt16(user.Gender),
		AvatarFileKey: user.AvatarFileKey.Ptr(),
	}
}

func ToClientDTO(client models.Client) dto.Client {
	return dto.Client{
		ID:        client.ID,
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}
}

func ToSessionDTO(session models.Session) dto.Session {
	return dto.Session{
		ID:             session.ID,
		UserID:         session.UserID,
		RefreshTokenID: session.RefreshToken,
		UserAgent:      session.UserAgent,
		Fingerprint:    session.Fingerprint,
		ExpiresAt:      session.ExpiresAt,
	}
}

func ToUserStorage(user entity.User, setters ...models.UserOption) models.User {
	userStorageModel := models.User{
		ID:            user.ID,
		Username:      user.Username,
		PasswordHash:  user.PasswordHash,
		FullName:      user.FullName,
		DateOfBirth:   null.TimeFromPtr(user.DateOfBirth),
		Gender:        user.Gender.ToNullInt16(),
		AvatarFileKey: null.StringFromPtr(user.AvatarFileKey),
	}

	for _, setter := range setters {
		setter(&userStorageModel)
	}

	return userStorageModel
}

func ToSessionStorage(session entity.Session, setters ...models.SessionOption) models.Session {
	sessionStorageModel := models.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshTokenID,
		UserAgent:    session.UserAgent,
		Fingerprint:  session.Fingerprint,
		ExpiresAt:    session.ExpiresAt,
	}

	for _, setter := range setters {
		setter(&sessionStorageModel)
	}

	return sessionStorageModel
}
