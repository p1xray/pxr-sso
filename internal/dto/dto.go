package dto

import (
	"pxr-sso/internal/domain"
	"time"
)

// UserDTO is information about the user.
type UserDTO struct {
	ID          int64
	Permissions []string
}

// ClientDTO is information about the client.
type ClientDTO struct {
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
	UserAgent    string
	Fingerprint  string
}

// TokensDTO represent auth tokens.
type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}
