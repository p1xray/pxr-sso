package domain

import "time"

type GenderEnum uint

// Gender enum.
const (
	MALE   GenderEnum = 1
	FEMALE GenderEnum = 2
)

// User is data for user in storage.
type User struct {
	ID            int64
	Username      string
	PasswordHash  string
	FIO           string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
	Deleted       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Client is data for client in storage.
type Client struct {
	ID        int64
	Name      string
	Code      string
	SecretKey string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session is data for session in storage.
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
