package dto

import (
	"fmt"
)

type AuthorizeContext struct {
	request  ValidatedAuthorizeRequest
	client   Client
	sessions []Session
}

func NewAuthorizeContext(request ValidatedAuthorizeRequest, client Client, sessions []Session) AuthorizeContext {
	return AuthorizeContext{
		request:  request,
		client:   client,
		sessions: sessions,
	}
}

func (a *AuthorizeContext) Request() ValidatedAuthorizeRequest {
	return a.request
}

func (a *AuthorizeContext) RequestPrompt() string {
	return a.request.Prompt()
}

func (a *AuthorizeContext) RequestGrantedScopes() []string {
	return a.request.GrantedScopes()
}

func (a *AuthorizeContext) RequestPendingScopes() []string {
	return a.request.PendingScopes()
}

func (a *AuthorizeContext) RequestRedirectURI() string {
	return a.request.RedirectURI()
}

func (a *AuthorizeContext) RequestState() string {
	return a.request.State()
}

func (a *AuthorizeContext) Client() Client {
	return a.client
}

func (a *AuthorizeContext) Sessions() []Session {
	return a.sessions
}

func (a *AuthorizeContext) SingleSession() (Session, error) {
	const op = "get single session"

	if len(a.sessions) == 0 {
		return Session{}, fmt.Errorf("%s: no sessions found in authorize context", op)
	}

	if len(a.sessions) > 1 {
		return Session{}, fmt.Errorf("%s: multiple sessions found in authorize context", op)
	}

	return a.sessions[0], nil
}
