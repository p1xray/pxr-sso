package dto

type ValidatedAuthorize struct {
	responseType        string
	clientID            string
	redirectURI         string
	codeChallenge       string
	codeChallengeMethod string
	state               string
	scope               []string
}

func NewValidatedAuthorize(
	responseType,
	clientID,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state string,
	scope []string,
) ValidatedAuthorize {
	return ValidatedAuthorize{
		responseType:        responseType,
		clientID:            clientID,
		redirectURI:         redirectURI,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		state:               state,
		scope:               scope,
	}
}

func (va *ValidatedAuthorize) ResponseType() string {
	return va.responseType
}

func (va *ValidatedAuthorize) ClientID() string {
	return va.clientID
}

func (va *ValidatedAuthorize) RedirectURI() string {
	return va.redirectURI
}

func (va *ValidatedAuthorize) CodeChallenge() string {
	return va.codeChallenge
}

func (va *ValidatedAuthorize) CodeChallengeMethod() string {
	return va.codeChallengeMethod
}

func (va *ValidatedAuthorize) State() string {
	return va.state
}

func (va *ValidatedAuthorize) Scope() []string {
	return va.scope
}
