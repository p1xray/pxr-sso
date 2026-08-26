package models

import (
	"time"
)

// Scope is data for scope in storage.
type Scope struct {
	ID          int64
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
