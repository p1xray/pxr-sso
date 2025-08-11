package register

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type Params struct {
	Username      string
	Password      string
	ClientCode    string
	FIO           string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
	UserAgent     string
	Fingerprint   string
	Issuer        string
}
