package models

import (
	"time"
)

// Role is data for role in storage.
type Role struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Active      bool
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	PermissionLinks []RolePermissionLink
}
