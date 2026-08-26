package postgresql

// Config is the PostgreSQL storage configuration.
type Config struct {
	ConnectionURL string `yaml:"connection_url" env:"PXR_SSO_POSTGRES_URL" env-required:"true" env-upd:""`
}
