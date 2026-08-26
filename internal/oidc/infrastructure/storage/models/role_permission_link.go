package models

import "time"

// RolePermissionLink is data for role permission link in storage.
type RolePermissionLink struct {
	ID           int64
	RoleID       int64
	PermissionID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Permission Permission
}
