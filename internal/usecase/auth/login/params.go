package login

// Params is a data for log in use-case.
type Params struct {
	Username    string
	Password    string
	ClientCode  string
	UserAgent   string
	Fingerprint string
	Issuer      string
}
