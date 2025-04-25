package dto

// UserDTO is information about the user.
type UserDTO struct {
	ID          int64
	Permissions []string
}

// ClientDTO is information about the client.
type ClientDTO struct {
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
	Issuer      string
}

// TokensDTO represent auth tokens.
type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}
