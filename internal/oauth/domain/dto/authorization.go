package dto

type Authorization struct {
	code          string
	username      string
	clientID      string
	redirectURI   string
	codeChallenge string
	scope         []string
}

func NewAuthorization(
	code,
	username,
	clientID,
	redirectURI,
	codeChallenge string,
	scope []string,
) Authorization {
	return Authorization{
		code:          code,
		username:      username,
		clientID:      clientID,
		redirectURI:   redirectURI,
		codeChallenge: codeChallenge,
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

func (a *Authorization) Scope() []string {
	return a.scope
}
