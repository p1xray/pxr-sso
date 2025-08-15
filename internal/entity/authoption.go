package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
)

type AuthOption func(*Auth) error

func WithUser(user dto.User) AuthOption {
	return func(a *Auth) error {
		if user.ID == emptyID {
			return nil
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

		return nil
	}
}

func WithClient(client dto.Client) AuthOption {
	return func(a *Auth) error {
		if client.ID == emptyID {
			return nil
		}

		a.client = client

		return nil
	}
}

func WithSession(sessions ...dto.Session) AuthOption {
	return func(a *Auth) error {
		sessionEntities := make([]Session, 0)
		for _, session := range sessions {
			if session.ID == emptyID {
				continue
			}

			sessionEntity, err := NewSession(
				session.UserID,
				session.UserAgent,
				session.Fingerprint,
				WithSessionID(session.ID),
				WithSessionRefreshTokenID(session.RefreshTokenID),
				WithSessionExpiresAt(session.ExpiresAt),
			)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCreateSession, err)
			}

			sessionEntities = append(sessionEntities, sessionEntity)
		}

		a.Sessions = append(a.Sessions, sessionEntities...)

		return nil
	}
}
