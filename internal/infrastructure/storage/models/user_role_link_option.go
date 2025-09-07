package models

import "time"

type UserRoleLinkOption func(*UserRoleLink)

func UserRoleLinkCreated() UserRoleLinkOption {
	now := time.Now()
	return func(url *UserRoleLink) {
		url.CreatedAt = now
		url.UpdatedAt = now
	}
}

func UserRoleLinkUpdated() UserRoleLinkOption {
	return func(url *UserRoleLink) {
		url.UpdatedAt = time.Now()
	}
}
