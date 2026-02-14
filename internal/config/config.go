package config

import (
	"time"
)

// Config is the project configuration.
type Config struct {
	Env      string         `yaml:"env" env:"ENV" env-default:"local" env-upd:""`
	GRPC     GRPCConfig     `yaml:"grpc" env-required:"true"`
	Postgres PostgresConfig `yaml:"postgres" env-required:"true"`
	Redis    RedisConfig    `yaml:"redis" env-required:"true"`
	Tokens   TokensConfig   `yaml:"tokens" env-required:"true"`
}

// GRPCConfig is the gRPC controller configuration.
type GRPCConfig struct {
	Port    string        `yaml:"port" evn:"PXR_SSO_GRPC_PORT" env-required:"true" env-upd:""`
	Timeout time.Duration `yaml:"timeout" env:"PXR_SSO_GRPC_TIMEOUT" env-required:"true" env-upd:""`
}

// PostgresConfig is the PostgreSQL storage configuration.
type PostgresConfig struct {
	ConnectionURL string `yaml:"connection_url" env:"PXR_SSO_POSTGRES_URL" env-required:"true" env-upd:""`
}

// RedisConfig is the redis storage configuration.
type RedisConfig struct {
	ConnectionURL string        `yaml:"connection_url" env:"PXR_SSO_REDIS_URL" env-required:"true" env-upd:""`
	FlowTTL       time.Duration `yaml:"flow_ttl" env:"PXR_SSO_REDIS_FLOW_TTL" env-required:"true" env-upd:""`
}

// TokensConfig is the auth tokens configuration.
type TokensConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"PXR_SSO_ACCESS_TOKEN_TTL" env-required:"true" env-upd:""`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"PXR_SSO_REFRESH_TOKEN_TTL" env-required:"true" env-upd:""`
}
