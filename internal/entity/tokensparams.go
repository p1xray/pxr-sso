package entity

import "time"

type CreateTokensParams struct {
	UserID          int64
	Permissions     []string
	ClientCode      string
	SecretKey       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
