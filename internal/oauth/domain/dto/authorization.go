package dto

type Authorization struct {
	code          string
	username      string
	clientID      string
	redirectURI   string
	codeChallenge string
	audience      string
	scope         []string
}

func NewAuthorization(
	code,
	username,
	clientID,
	redirectURI,
	codeChallenge,
	audience string,
	scope []string,
) Authorization {
	return Authorization{
		code:          code,
		username:      username,
		clientID:      clientID,
		redirectURI:   redirectURI,
		codeChallenge: codeChallenge,
		audience:      audience,
		scope:         scope,
	}
}

func (a *Authorization) Code() string {
	return a.code
}

func (a *Authorization) Username() string {
	return a.username
}

func (a *Authorization) ClientID() string {
	return a.clientID
}

func (a *Authorization) RedirectURI() string {
	return a.redirectURI
}

func (a *Authorization) CodeChallenge() string {
	return a.codeChallenge
}

func (a *Authorization) Audience() string {
	return a.audience
}

func (a *Authorization) Scope() []string {
	return a.scope
}
