package dto

import (
	"pxr-sso/internal/domain"
	"time"
)

// UserDTO is information about the user.
type UserDTO struct {
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
	Gender        *domain.GenderEnum
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
	Gender        *domain.GenderEnum
	AvatarFileKey *string
}
