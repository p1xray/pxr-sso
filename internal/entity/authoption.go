package entity

import "github.com/p1xray/pxr-sso/internal/dto"

const emptyID = 0

type AuthOption func(*Auth)

func WithUser(user dto.User) AuthOption {
	return func(a *Auth) {
		if user.ID == emptyID {
			return
		}

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
		if client.ID == emptyID {
			return
		}

		a.client = client
	}
}

func WithSession(sessions ...dto.Session) AuthOption {
	return func(a *Auth) {
		sessionEntities := make([]Session, 0)
		for _, session := range sessions {
			if session.ID == emptyID {
				continue
			}

			sessionEntity := NewExistSession(
				session.ID,
				session.UserID,
				session.RefreshTokenID,
				session.UserAgent,
				session.Fingerprint,
				session.ExpiresAt,
			)

			sessionEntities = append(sessionEntities, sessionEntity)
		}

		a.Sessions = append(a.Sessions, sessionEntities...)
	}
}
