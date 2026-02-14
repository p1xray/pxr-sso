package app

import (
	grpcapp "github.com/p1xray/pxr-sso/internal/app/grpc"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/postgresql"
)

// Config is the project configuration.
type Config struct {
	Env        string                   `yaml:"env" env:"ENV" env-default:"local" env-upd:""`
	GRPC       grpcapp.Config           `yaml:"grpc" env-required:"true"`
	Postgres   postgresql.Config        `yaml:"postgres" env-required:"true"`
	Redis      redis.Config             `yaml:"redis" env-required:"true"`
	Token      generator.TokenConfig    `yaml:"token" env-required:"true"`
	URIBuilder builder.URIBuilderConfig `yaml:"uri_builder" env-required:"true"`
}
