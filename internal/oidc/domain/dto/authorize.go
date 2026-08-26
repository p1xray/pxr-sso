package dto

import (
	"github.com/p1xray/pxr-sso/internal/oidc"
	"strings"
)

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
	sessions            []SessionCookie
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
	scope []string,
	sessions []SessionCookie,
) AuthorizeRequest {
	if len(prompt) == 0 {
		prompt = []string{oidc.DefaultPrompt}
	}

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
		sessions:            sessions,
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
	scopes := make([]string, 0)
	for _, scope := range a.scope {
		scopes = append(scopes, strings.Split(scope, " ")...)
	}

	return scopes
}

func (a *AuthorizeRequest) Sessions() []SessionCookie {
	return a.sessions
}
