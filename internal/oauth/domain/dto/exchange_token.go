package dto

type ExchangeToken struct {
	flowID            string
	grantType         string
	clientID          string
	authorizationCode string
	redirectURI       string
	codeVerifier      string
	audience          string
	scope             []string
}

func NewExchangeToken(
	flowID,
	grantType,
	clientID,
	authorizationCode,
	redirectURI,
	codeVerifier,
	audience string,
	scope []string,
) ExchangeToken {
	return ExchangeToken{
		flowID:            flowID,
		grantType:         grantType,
		clientID:          clientID,
		authorizationCode: authorizationCode,
		redirectURI:       redirectURI,
		codeVerifier:      codeVerifier,
		audience:          audience,
		scope:             scope,
	}
}

func (t *ExchangeToken) FlowID() string {
	return t.flowID
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

func (t *ExchangeToken) Audience() string {
	return t.audience
}

func (t *ExchangeToken) Scope() []string {
	return t.scope
}
