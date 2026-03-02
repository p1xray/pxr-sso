package token

type Params struct {
	GrantType         string
	ClientID          string
	ClientSecret      string
	AuthorizationCode string
	RedirectURI       string
	CodeVerifier      string
	Audience          string
	Scope             []string
}
