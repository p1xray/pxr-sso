package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           int64
	username     string
	passwordHash string
	fullName     string
	roles        []dto.Role
}

func NewUser(username, password, fullName string, defaultRoles []dto.Role) (User, error) {
	passwordHash, err := generator.PasswordHash(password)
	if err != nil {
		return User{}, fmt.Errorf("generate password hash: %w", err)
	}

	return User{
		id:           0,
		username:     username,
		passwordHash: passwordHash,
		fullName:     fullName,
		roles:        defaultRoles,
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

func (u *User) CheckCredentials(username, password string) error {
	if u.Username() != username {
		return oidc.ErrOAuthInvalidUserCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash()), []byte(password)); err != nil {
		return oidc.ErrOAuthInvalidUserCredentials
	}

	return nil
}
