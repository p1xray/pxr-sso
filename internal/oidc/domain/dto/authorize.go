package dto

type AuthorizeRequest struct {
	responseType        []string
	prompt              []string
	clientID            []string
	redirectURI         []string
	codeChallenge       []string
	codeChallengeMethod []string
	state               []string
	audience            []string
	scope               []string
	session             []string
}

func NewAuthorizeRequest(
	responseType,
	prompt,
	clientID,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state,
	audience,
	scope,
	session []string,
) AuthorizeRequest {
	return AuthorizeRequest{
		responseType:        responseType,
		prompt:              prompt,
		clientID:            clientID,
		redirectURI:         redirectURI,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		state:               state,
		audience:            audience,
		scope:               scope,
		session:             session,
	}
}

func (a *AuthorizeRequest) ResponseType() []string {
	return a.responseType
}

func (a *AuthorizeRequest) Prompt() []string {
	return a.prompt
}

func (a *AuthorizeRequest) ClientID() []string {
	return a.clientID
}

func (a *AuthorizeRequest) RedirectURI() []string {
	return a.redirectURI
}

func (a *AuthorizeRequest) CodeChallenge() []string {
	return a.codeChallenge
}

func (a *AuthorizeRequest) CodeChallengeMethod() []string {
	return a.codeChallengeMethod
}

func (a *AuthorizeRequest) State() []string {
	return a.state
}

func (a *AuthorizeRequest) Audience() []string {
	return a.audience
}

func (a *AuthorizeRequest) Scope() []string {
	return a.scope
}

func (a *AuthorizeRequest) Session() []string {
	return a.session
}
