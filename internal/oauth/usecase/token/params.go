package token

type Params struct {
	FlowID            string
	GrantType         string
	ClientID          string
	AuthorizationCode string
	RedirectURI       string
	CodeVerifier      string
	Scope             []string
}
