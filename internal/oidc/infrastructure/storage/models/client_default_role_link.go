package models

import "time"

// ClientDefaultRoleLink is data for client default role link in storage.
type ClientDefaultRoleLink struct {
	ID        int64
	ClientID  int64
	RoleID    int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role
}
