package models

import "time"

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
