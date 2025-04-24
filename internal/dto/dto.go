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

// LoginDTO is data for login user.
type LoginDTO struct {
	Username    string
	Password    string
	ClientCode  string
	UserAgent   string
	Fingerprint string
}

// TokensDTO represent auth tokens.
type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}
