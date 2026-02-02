package dto

type Login struct {
	flowID       string
	responseType string
	clientID     string
	redirectURI  string
	state        string
	username     string
	password     string
}

func NewLogin(
	flowID,
	responseType,
	clientID,
	redirectURI,
	state,
	username,
	password string,
) Login {
	return Login{
		flowID:       flowID,
		responseType: responseType,
		clientID:     clientID,
		redirectURI:  redirectURI,
		state:        state,
		username:     username,
		password:     password,
	}
}

func (l Login) FlowID() string {
	return l.flowID
}

func (l Login) ResponseType() string {
	return l.responseType
}

func (l Login) ClientID() string {
	return l.clientID
}

func (l Login) RedirectURI() string {
	return l.redirectURI
}

func (l Login) State() string {
	return l.state
}

func (l Login) Username() string {
	return l.username
}

func (l Login) Password() string {
	return l.password
}
