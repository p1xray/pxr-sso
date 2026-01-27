package edit

import (
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

// Params is a data for edit user profile data use-case.
type Params struct {
	ID            int64
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
}
