package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
)

// AuthOption is how options for the Auth are set up.
type AuthOption func(*Auth) error

// WithAuthUser is an option which sets up the user data for the user authentication entity.
func WithAuthUser(user dto.User) AuthOption {
	return func(a *Auth) error {
		if user.ID == emptyID {
			return nil
		}

		a.User = NewUser(
			user.Username,
			user.FullName,
			user.DateOfBirth,
			user.Gender,
			user.AvatarFileKey,
			WithUserID(user.ID),
			WithUserPasswordHash(user.PasswordHash),
			WithUserRoles(user.Roles),
			WithUserPermissions(user.Permissions),
		)

		return nil
	}
}

// WithAuthClient is an option which sets up the client data for the user authentication entity.
func WithAuthClient(client dto.Client) AuthOption {
	return func(a *Auth) error {
		if client.ID == emptyID {
			return nil
		}

		a.client = client

		return nil
	}
}

// WithAuthSession is an option which sets up the sessions data for the user authentication entity.
func WithAuthSession(sessions ...dto.Session) AuthOption {
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

// WithAuthDefaultRoles is an option which sets up the default roles for the user authentication entity.
func WithAuthDefaultRoles(roles ...dto.Role) AuthOption {
	return func(a *Auth) error {
		a.defaultRoles = roles

		return nil
	}
}

// WithAuthDefaultPermissionCodes is an option which sets up the default permission codes
// for the user authentication entity.
func WithAuthDefaultPermissionCodes(permissions ...string) AuthOption {
	return func(a *Auth) error {
		a.defaultPermissionCodes = permissions

		return nil
	}
}
