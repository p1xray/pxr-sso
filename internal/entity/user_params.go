package entity

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type UserUpdateProfileParams struct {
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
}
