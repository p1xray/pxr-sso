package entity

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Permissions  []string
}

func NewUser(id int64, username, passwordHash string, permissions []string) User {
	return User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		Permissions:  permissions,
	}
}
