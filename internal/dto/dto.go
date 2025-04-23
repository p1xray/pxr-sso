package dto

// User is information about the user.
type User struct {
	ID          int64
	Permissions []string
}

// Client is information about the client.
type Client struct {
	Code      string
	SecretKey string
}
