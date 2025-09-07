package dto

// Client is a DTO with client data.
type Client struct {
	ID        int64
	Code      string
	SecretKey string
	Audiences []string
}
