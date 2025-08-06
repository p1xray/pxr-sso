package entity

import "time"

type CreateNewSessionParams struct {
	UserID          int64
	UserPermissions []string
	ClientCode      string
	ClientSecretKey string
	Issuer          string
	UserAgent       string
	Fingerprint     string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
