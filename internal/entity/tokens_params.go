package entity

import "time"

// CreateTokensParams is a data for creating new user session tokens.
type CreateTokensParams struct {
	UserID          int64
	Permissions     []string
	Audiences       []string
	SecretKey       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
