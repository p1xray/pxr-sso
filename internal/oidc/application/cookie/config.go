package cookie

// SessionConfig is the configuration for session cookie.
type SessionConfig struct {
	SecretKey string `yaml:"secret_key" env:"PXR_SSO_SESSION_COOKIE_SECRET_KEY" env-required:"true" env-upd:""`
}
