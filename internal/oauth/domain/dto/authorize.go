package dto

type Authorize struct {
	responseType        []string
	clientID            []string
	redirectURI         []string
	codeChallenge       []string
	codeChallengeMethod []string
	state               []string
}

func NewAuthorize(
	responseType,
	clientID,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state []string,
) Authorize {
	return Authorize{
		responseType:        responseType,
		clientID:            clientID,
		redirectURI:         redirectURI,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		state:               state,
	}
}

func (a *Authorize) ResponseType() []string {
	return a.responseType
}

func (a *Authorize) ClientID() []string {
	return a.clientID
}

func (a *Authorize) RedirectURI() []string {
	return a.redirectURI
}

func (a *Authorize) CodeChallenge() []string {
	return a.codeChallenge
}

func (a *Authorize) CodeChallengeMethod() []string {
	return a.codeChallengeMethod
}

func (a *Authorize) State() []string {
	return a.state
}
