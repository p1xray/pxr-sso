package entity

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

// LoginParams is a data for logging in a user.
type LoginParams struct {
	Password    string
	UserAgent   string
	Fingerprint string
	Issuer      string
}

// RegisterParams is a data for registering a user.
type RegisterParams struct {
	Username      string
	Password      string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
	UserAgent     string
	Fingerprint   string
	Issuer        string
}

// RefreshTokensParams is a data for refreshing user tokens.
type RefreshTokensParams struct {
	UserAgent   string
	Fingerprint string
	Issuer      string
}
