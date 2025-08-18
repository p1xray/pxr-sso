package dto

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

// User is information about the user.
type User struct {
	ID            int64
	Username      string
	PasswordHash  string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
	Roles         []Role
	Permissions   []string
}

// UserProfile is data for user profile.
type UserProfile struct {
	ID            int64
	Username      string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
}
