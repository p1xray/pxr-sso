package entity

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type LoginParams struct {
	Password    string
	UserAgent   string
	Fingerprint string
	Issuer      string
}

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

type RefreshTokensParams struct {
	UserAgent   string
	Fingerprint string
	Issuer      string
}
