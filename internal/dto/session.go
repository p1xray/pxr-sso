package dto

import "time"

type Session struct {
	ID             int64
	UserID         int64
	RefreshTokenID string
	UserAgent      string
	Fingerprint    string
	ExpiresAt      time.Time
}
