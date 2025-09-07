package entity

import "github.com/p1xray/pxr-sso/internal/dto"

// UserOption is how options for the User are set up.
type UserOption func(*User)

// WithUserID is an option which sets up the user ID for the user entity.
func WithUserID(id int64) UserOption {
	return func(u *User) {
		u.ID = id
	}
}

// WithUserPasswordHash is an option which sets up the user password hash for the user entity.
func WithUserPasswordHash(passwordHash string) UserOption {
	return func(u *User) {
		u.PasswordHash = passwordHash
	}
}

// WithUserRoles is an option which sets up the user roles for the user entity.
func WithUserRoles(roles []dto.Role) UserOption {
	return func(u *User) {
		u.Roles = roles
	}
}

// WithUserPermissions is an option which sets up the user permissions for the user entity.
func WithUserPermissions(permissions []string) UserOption {
	return func(u *User) {
		u.Permissions = permissions
	}
}
