package dto

type Client struct {
	ID        int64
	Code      string
	SecretKey string
	Audiences []string
}
