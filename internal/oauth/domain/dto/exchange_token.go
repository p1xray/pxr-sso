package dto

type ExchangeToken struct {
	grantType         string
	clientID          string
	clientSecret      string
	authorizationCode string
	redirectURI       string
	codeVerifier      string
	audience          string
	scope             []string
}

func NewExchangeToken(
	grantType,
	clientID,
	clientSecret,
	authorizationCode,
	redirectURI,
	codeVerifier,
	audience string,
	scope []string,
) ExchangeToken {
	return ExchangeToken{
		grantType:         grantType,
		clientID:          clientID,
		clientSecret:      clientSecret,
		authorizationCode: authorizationCode,
		redirectURI:       redirectURI,
		codeVerifier:      codeVerifier,
		audience:          audience,
		scope:             scope,
	}
}

func (t *ExchangeToken) GrantType() string {
	return t.grantType
}

func (t *ExchangeToken) ClientID() string {
	return t.clientID
}

func (t *ExchangeToken) ClientSecret() string {
	return t.clientSecret
}

func (t *ExchangeToken) AuthorizationCode() string {
	return t.authorizationCode
}

func (t *ExchangeToken) RedirectURI() string {
	return t.redirectURI
}

func (t *ExchangeToken) CodeVerifier() string {
	return t.codeVerifier
}

func (t *ExchangeToken) Audience() string {
	return t.audience
}

func (t *ExchangeToken) Scope() []string {
	return t.scope
}
