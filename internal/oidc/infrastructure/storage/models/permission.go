package models

import "time"

// Permission is data for permission in storage.
type Permission struct {
	ID          int64
	Code        string
	Description string
	Active      bool
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
