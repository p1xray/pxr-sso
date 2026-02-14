package generator

import "time"

// TokenConfig is the auth tokens configuration.
type TokenConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"PXR_SSO_TOKEN_ACCESS_TOKEN_TTL" env-required:"true" env-upd:""`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"PXR_SSO_TOKEN_REFRESH_TOKEN_TTL" env-required:"true" env-upd:""`
	Issuer          string        `yaml:"issuer" env:"PXR_SSO_TOKEN_ISSUER" env-required:"true" env-upd:""`
}
