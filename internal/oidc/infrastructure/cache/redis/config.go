package redis

import "time"

// Config is the redis storage configuration.
type Config struct {
	ConnectionURL           string        `yaml:"connection_url" env:"PXR_SSO_REDIS_URL" env-required:"true" env-upd:""`
	AuthorizationRequestTTL time.Duration `yaml:"authorization_request_ttl" env:"PXR_SSO_REDIS_AUTHORIZATION_REQUEST_TTL" env-required:"true" env-upd:""`
	AuthorizedGrantTTL      time.Duration `yaml:"authorized_grant_ttl" env:"PXR_SSO_REDIS_AUTHORIZED_GRANT_TTL" env-required:"true" env-upd:""`
}
