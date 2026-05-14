package redis

import "time"

// Config is the redis storage configuration.
type Config struct {
	ConnectionURL    string        `yaml:"connection_url" env:"PXR_SSO_REDIS_URL" env-required:"true" env-upd:""`
	FlowTTL          time.Duration `yaml:"flow_ttl" env:"PXR_SSO_REDIS_FLOW_TTL" env-required:"true" env-upd:""`
	AuthorizationTTL time.Duration `yaml:"authorization_ttl" env:"PXR_SSO_REDIS_AUTHORIZATION_TTL" env-required:"true" env-upd:""`
}
