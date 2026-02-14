package grpcapp

import "time"

// Config is the gRPC controller configuration.
type Config struct {
	Port    string        `yaml:"port" evn:"PXR_SSO_GRPC_PORT" env-required:"true" env-upd:""`
	Timeout time.Duration `yaml:"timeout" env:"PXR_SSO_GRPC_TIMEOUT" env-required:"true" env-upd:""`
}
