package dto

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

// User is a DTO with user data.
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

// UserProfile is a DTO with user profile data.
type UserProfile struct {
	ID            int64
	Username      string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
}
