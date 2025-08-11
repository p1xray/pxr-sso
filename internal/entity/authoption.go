package entity

import "github.com/p1xray/pxr-sso/internal/dto"

type AuthOption func(*Auth)

func WithUser(user dto.User) AuthOption {
	return func(a *Auth) {
		a.User = NewUser(
			user.ID,
			user.Username,
			user.PasswordHash,
			user.FullName,
			user.DateOfBirth,
			user.Gender,
			user.AvatarFileKey,
			user.Roles,
			user.Permissions,
		)
	}
}

func WithClient(client dto.Client) AuthOption {
	return func(a *Auth) {
		a.client = client
	}
}

func WithSession(sessions ...dto.Session) AuthOption {
	return func(a *Auth) {
		sessionEntities := make([]Session, len(sessions))
		for i, session := range sessions {
			sessionEntity := NewExistSession(
				session.ID,
				session.UserID,
				session.RefreshTokenID,
				session.UserAgent,
				session.Fingerprint,
				session.ExpiresAt,
			)

			sessionEntities[i] = sessionEntity
		}

		a.Sessions = append(a.Sessions, sessionEntities...)
	}
}
