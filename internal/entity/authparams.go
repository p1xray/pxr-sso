package entity

import (
	"github.com/p1xray/pxr-sso/internal/dto"
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
	ClientCode    string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *dto.GenderEnum
	AvatarFileKey *string
	UserAgent     string
	Fingerprint   string
	Issuer        string
}
