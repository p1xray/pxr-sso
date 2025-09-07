package models

import "time"

type Role struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Active      bool
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
