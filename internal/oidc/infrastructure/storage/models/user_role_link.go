package models

import "time"

// UserRoleLink is data for user role link in storage.
type UserRoleLink struct {
	ID        int64
	UserID    int64
	RoleID    int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role
}

type UserRoleLinkOption func(*UserRoleLink)

func UserRoleLinkCreated() UserRoleLinkOption {
	now := time.Now()
	return func(url *UserRoleLink) {
		url.CreatedAt = now
		url.UpdatedAt = now
	}
}
