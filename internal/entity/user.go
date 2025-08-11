package entity

import (
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
	Roles         []string
	Permissions   []string

	dataStatus enum.DataStatusEnum
}

func NewUser(
	id int64,
	username,
	passwordHash,
	fullName string,
	dateOfBirth *time.Time,
	gender *enum.GenderEnum,
	avatarFileKey *string,
	roles []string,
	permissions []string,
) User {
	return User{
		ID:            id,
		Username:      username,
		PasswordHash:  passwordHash,
		FullName:      fullName,
		DateOfBirth:   dateOfBirth,
		Gender:        gender,
		AvatarFileKey: avatarFileKey,
		Roles:         roles,
		Permissions:   permissions,
	}
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
