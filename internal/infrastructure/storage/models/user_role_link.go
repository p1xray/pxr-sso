package models

import "time"

type UserRoleLink struct {
	ID        int64
	UserID    int64
	RoleID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
