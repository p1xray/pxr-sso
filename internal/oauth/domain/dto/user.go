package dto

import "strconv"

// User is a DTO with user data.
type User struct {
	id           int64
	username     string
	passwordHash string
	fullName     string
	roles        []Role
}

func NewUser(id int64, username, passwordHash string) User {
	return User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
	}
}

func NewRegisteringUser(username, passwordHash, fullName string, roles []Role) User {
	return User{
		username:     username,
		passwordHash: passwordHash,
		fullName:     fullName,
		roles:        roles,
	}
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) IDString() string {
	return strconv.FormatInt(u.id, 10)
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

func (u *User) Roles() []Role {
	return u.roles
}
