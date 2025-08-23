package entity

import "time"

type SessionWithGeneratedTokensParams struct {
	UserPermissions []string
	ClientCode      string
	ClientSecretKey string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
