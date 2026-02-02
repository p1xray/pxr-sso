package login

import "github.com/google/uuid"

type Params struct {
	FlowID       uuid.UUID
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
	Scope        string
	Username     string
	Password     string
}
