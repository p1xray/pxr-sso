package dto

import (
	"github.com/google/uuid"
)

type Flow struct {
	id                  uuid.UUID
	clientID            string
	responseType        string
	redirectURI         string
	codeChallenge       string
	codeChallengeMethod string
	state               string
	username            string
	scope               []string
}

func NewFlow(
	id uuid.UUID,
	clientID,
	responseType,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state string,
	scope []string,
	setters ...FlowOption,
) Flow {
	flow := Flow{
		id:                  id,
		clientID:            clientID,
		responseType:        responseType,
		redirectURI:         redirectURI,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		state:               state,
		scope:               scope,
	}

	for _, setter := range setters {
		setter(&flow)
	}

	return flow
}

func (f *Flow) ID() uuid.UUID {
	return f.id
}

func (f *Flow) ClientID() string {
	return f.clientID
}

func (f *Flow) ResponseType() string {
	return f.responseType
}

func (f *Flow) RedirectURI() string {
	return f.redirectURI
}

func (f *Flow) CodeChallenge() string {
	return f.codeChallenge
}

func (f *Flow) CodeChallengeMethod() string {
	return f.codeChallengeMethod
}

func (f *Flow) State() string {
	return f.state
}

func (f *Flow) Username() string {
	return f.username
}

func (f *Flow) Scope() []string {
	return f.scope
}

// FlowOption is how options for the Flow are set up.
type FlowOption func(*Flow)

// WithUsername is an option which sets up the username for the Flow.
func WithUsername(username string) FlowOption {
	return func(a *Flow) {
		a.username = username
	}
}
