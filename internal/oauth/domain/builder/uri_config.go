package builder

// URIBuilderConfig is the URI builder configuration.
type URIBuilderConfig struct {
	Login   string `yaml:"login" env:"PXR_SSO_URI_BUILDER_LOGIN" env-required:"true" env-upd:""`
	Consent string `yaml:"consent" env:"PXR_SSO_URI_BUILDER_CONSENT" env-required:"true" env-upd:""`
	Error   string `yaml:"error" env:"PXR_SSO_URI_BUILDER_ERROR" env-required:"true" env-upd:""`
}
