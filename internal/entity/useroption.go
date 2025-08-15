package entity

type UserOption func(*User)

func WithUserID(id int64) UserOption {
	return func(u *User) {
		u.ID = id
	}
}

func WithUserPasswordHash(passwordHash string) UserOption {
	return func(u *User) {
		u.PasswordHash = passwordHash
	}
}

func WithUserRoles(roles []string) UserOption {
	return func(u *User) {
		u.Roles = roles
	}
}

func WithUserPermissions(permissions []string) UserOption {
	return func(u *User) {
		u.Permissions = permissions
	}
}
