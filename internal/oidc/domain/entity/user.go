package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/enum"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	id            int64
	username      string
	passwordHash  string
	fullName      string
	dateOfBirth   nullable.Nullable[time.Time]
	gender        nullable.Nullable[enum.Gender]
	avatarFileKey nullable.Nullable[string]
	clientID      int64
	roles         []dto.Role
	dataStatus    enum.DataStatus
}

func NewUser(
	username,
	password,
	fullName string,
	clientID int64,
	defaultRoles []dto.Role,
) (User, error) {
	passwordHash, err := generator.PasswordHash(password)
	if err != nil {
		return User{}, fmt.Errorf("generate password hash: %w", err)
	}

	return User{
		id:            0,
		username:      username,
		passwordHash:  passwordHash,
		fullName:      fullName,
		dateOfBirth:   nullable.None[time.Time](),
		gender:        nullable.None[enum.Gender](),
		avatarFileKey: nullable.None[string](),
		clientID:      clientID,
		roles:         defaultRoles,
		dataStatus:    enum.DataStatusToCreate,
	}, nil
}

func NewExistUser(data dto.User) (User, error) {
	if data.ID() == 0 {
		return User{}, fmt.Errorf("user id is empty")
	}

	return User{
		id:           data.ID(),
		username:     data.Username(),
		passwordHash: data.PasswordHash(),
		fullName:     data.FullName(),
		// dateOfBirth:   data.DateOfBirth(),   // TODO: implement method
		// gender:        data.Gender(),        // TODO: implement method
		// avatarFileKey: data.AvatarFileKey(), // TODO: implement method
	}, nil
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) FullName() string {
	return u.fullName
}

func (u *User) DateOfBirth() nullable.Nullable[time.Time] {
	return u.dateOfBirth
}

func (u *User) Gender() nullable.Nullable[enum.Gender] {
	return u.gender
}

func (u *User) AvatarFileKey() nullable.Nullable[string] {
	return u.avatarFileKey
}

func (u *User) ClientID() int64 {
	return u.clientID
}

func (u *User) Roles() []dto.Role {
	return u.roles
}

func (u *User) DataStatus() enum.DataStatus {
	return u.dataStatus
}

func (u *User) CheckCredentials(username, password string) error {
	if u.Username() != username {
		return oidc.ErrOAuthInvalidUserCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash()), []byte(password)); err != nil {
		return oidc.ErrOAuthInvalidUserCredentials
	}

	return nil
}
