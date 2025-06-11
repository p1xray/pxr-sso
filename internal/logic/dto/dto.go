package dto

import (
	"github.com/guregu/null/v6"
	"time"
)

// GenderEnum is type for gender enum.
type GenderEnum int16

// Gender enum.
const (
	MALE   GenderEnum = 1
	FEMALE GenderEnum = 2
)

func (ge *GenderEnum) ToNullInt16() null.Int16 {
	if ge == nil {
		return null.NewInt16(0, false)
	}
	return null.Int16From(int16(*ge))
}

// UserWithPermissionsDTO is information about the user with permissions.
type UserWithPermissionsDTO struct {
	ID           int64
	PasswordHash string
	Permissions  []string
}

// ClientDTO is information about the client.
type ClientDTO struct {
	ID        int64
	Code      string
	SecretKey string
}

// LoginDTO is data for login user.
type LoginDTO struct {
	Username    string
	Password    string
	ClientCode  string
	UserAgent   string
	Fingerprint string
	Issuer      string
}

// RegisterDTO is data for register new user.
type RegisterDTO struct {
	Username      string
	Password      string
	ClientCode    string
	FIO           string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
	UserAgent     string
	Fingerprint   string
	Issuer        string
}

// RefreshTokensDTO is data for refresh user's auth tokens.
type RefreshTokensDTO struct {
	RefreshToken string
	ClientCode   string
	UserAgent    string
	Fingerprint  string
	Issuer       string
}

// LogoutDTO is data for logout.
type LogoutDTO struct {
	RefreshToken string
	ClientCode   string
}

// TokensDTO represent auth tokens.
type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}

// SessionDTO is information about the session.
type SessionDTO struct {
	ID           int64
	UserID       int64
	RefreshToken string
	UserAgent    string
	Fingerprint  string
	ExpiresAt    time.Time
}

// CreateSessionDTO is data for create new session.
type CreateSessionDTO struct {
	UserID       int64
	RefreshToken string
	UserAgent    string
	Fingerprint  string
	ExpiresAt    time.Time
}

// CreateUserDTO is data for create new user.
type CreateUserDTO struct {
	Username      string
	PasswordHash  []byte
	FIO           string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
}

// UserProfileDTO is user profile data.
type UserProfileDTO struct {
	UserID        int64
	Username      string
	FIO           string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
}

// UserDTO is information about the user.
type UserDTO struct {
	ID            int64
	Username      string
	FIO           string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
}
