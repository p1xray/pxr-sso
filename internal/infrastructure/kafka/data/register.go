package data

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type RegisteredUser struct {
	ID            int64            `json:"id"`
	FullName      string           `json:"full_name"`
	DateOfBirth   *time.Time       `json:"date_of_birth"`
	Gender        *enum.GenderEnum `json:"gender"`
	AvatarFileKey *string          `json:"avatar_file_key"`
}
