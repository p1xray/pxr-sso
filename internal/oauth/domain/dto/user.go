package dto

// User is a DTO with user data.
type User struct {
	id           int64
	username     string
	passwordHash string
}

func NewUser(id int64, username, passwordHash string) User {
	return User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
	}
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
