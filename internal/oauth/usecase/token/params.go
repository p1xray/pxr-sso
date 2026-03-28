package token

type Params struct {
	GrantType         string
	ClientID          string
	AuthorizationCode string
	RedirectURI       string
	CodeVerifier      string
}
