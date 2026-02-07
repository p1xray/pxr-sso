package models

import (
	"github.com/guregu/null/v6"
	"time"
)

// User is data for user in storage.
type User struct {
	ID            int64
	Username      string
	PasswordHash  string
	FullName      string
	DateOfBirth   null.Time
	Gender        null.Int16
	AvatarFileKey null.String
	Deleted       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserOption func(*User)

func UserCreated() UserOption {
	now := time.Now()
	return func(a *User) {
		a.Deleted = false
		a.CreatedAt = now
		a.UpdatedAt = now
	}
}

func UserUpdated() UserOption {
	return func(a *User) {
		a.Deleted = false
		a.UpdatedAt = time.Now()
	}
}

func UserRemoved() UserOption {
	return func(a *User) {
		a.Deleted = true
		a.UpdatedAt = time.Now()
	}
}
