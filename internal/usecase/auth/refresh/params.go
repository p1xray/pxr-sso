package refresh

// Params is a data for refresh user tokens use-case.
type Params struct {
	RefreshToken string
	ClientCode   string
	UserAgent    string
	Fingerprint  string
	Issuer       string
}
