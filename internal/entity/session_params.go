package entity

import "time"

// SessionWithGeneratedTokensParams is a data for option which sets up the generated tokens for the user session entity.
type SessionWithGeneratedTokensParams struct {
	UserPermissions []string
	Audiences       []string
	ClientSecretKey string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
