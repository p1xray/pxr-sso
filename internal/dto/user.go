package dto

// User is information about the user.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Permissions  []string
	Client       Client
}

// UserCredentials is user's credentials data.
type UserCredentials struct {
	ID           int64
	Username     string
	PasswordHash string
}
