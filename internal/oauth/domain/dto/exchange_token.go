package dto

type ExchangeToken struct {
	grantType         string
	clientID          string
	authorizationCode string
	redirectURI       string
	codeVerifier      string
}

func NewExchangeToken(
	grantType,
	clientID,
	authorizationCode,
	redirectURI,
	codeVerifier string,
) ExchangeToken {
	return ExchangeToken{
		grantType:         grantType,
		clientID:          clientID,
		authorizationCode: authorizationCode,
		redirectURI:       redirectURI,
		codeVerifier:      codeVerifier,
	}
}

func (t *ExchangeToken) GrantType() string {
	return t.grantType
}

func (t *ExchangeToken) ClientID() string {
	return t.clientID
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
