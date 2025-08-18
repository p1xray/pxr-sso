package entity

import (
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type User struct {
	ID            int64
	Username      string
	PasswordHash  string
	FullName      string
	DateOfBirth   *time.Time
	Gender        *enum.GenderEnum
	AvatarFileKey *string
	Roles         []dto.Role
	Permissions   []string

	dataStatus enum.DataStatusEnum
}

func NewUser(
	username,
	fullName string,
	dateOfBirth *time.Time,
	gender *enum.GenderEnum,
	avatarFileKey *string,
	setters ...UserOption,
) User {
	user := User{
		Username:      username,
		FullName:      fullName,
		DateOfBirth:   dateOfBirth,
		Gender:        gender,
		AvatarFileKey: avatarFileKey,
	}

	for _, setter := range setters {
		setter(&user)
	}

	return user
}

func (u *User) SetToCreate() {
	u.dataStatus = enum.ToCreate
}

func (u *User) SetToUpdate() {
	u.dataStatus = enum.ToUpdate
}

func (u *User) SetToRemove() {
	u.dataStatus = enum.ToRemove
}

func (u *User) IsToCreate() bool {
	return u.dataStatus == enum.ToCreate
}

func (u *User) IsToUpdate() bool {
	return u.dataStatus == enum.ToUpdate
}

func (u *User) IsToRemove() bool {
	return u.dataStatus == enum.ToRemove
}
