package register

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

// Params is a data for register use-case.
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
