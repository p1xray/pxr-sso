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
	authorizationCode   string
	userID              int64
}

func NewFlow(
	id uuid.UUID,
	clientID,
	responseType,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state string,
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

func (f *Flow) AuthorizationCode() string {
	return f.authorizationCode
}

func (f *Flow) UserID() int64 {
	return f.userID
}

// FlowOption is how options for the Flow are set up.
type FlowOption func(*Flow)

// WithAuthorizationCode is an option which sets up the authorization code for the Flow.
func WithAuthorizationCode(authorizationCode string) FlowOption {
	return func(a *Flow) {
		a.authorizationCode = authorizationCode
	}
}

// WithUserID is an option which sets up the user ID for the Flow.
func WithUserID(userID int64) FlowOption {
	return func(a *Flow) {
		a.userID = userID
	}
}
