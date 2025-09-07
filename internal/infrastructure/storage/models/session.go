package models

import "time"

type Session struct {
	ID           int64
	UserID       int64
	RefreshToken string
	UserAgent    string
	Fingerprint  string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
