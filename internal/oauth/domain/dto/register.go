package dto

type Register struct {
	flowID       string
	responseType string
	clientID     string
	redirectURI  string
	state        string
	username     string
	password     string
	fullName     string
	scope        []string
}

func NewRegister(
	flowID,
	responseType,
	clientID,
	redirectURI,
	state,
	username,
	password,
	fullName string,
	scope []string,
) Register {
	return Register{
		flowID:       flowID,
		responseType: responseType,
		clientID:     clientID,
		redirectURI:  redirectURI,
		state:        state,
		username:     username,
		password:     password,
		fullName:     fullName,
		scope:        scope,
	}
}

func (r Register) FlowID() string {
	return r.flowID
}

func (r Register) ResponseType() string {
	return r.responseType
}

func (r Register) ClientID() string {
	return r.clientID
}

func (r Register) RedirectURI() string {
	return r.redirectURI
}

func (r Register) State() string {
	return r.state
}

func (r Register) Scope() []string {
	return r.scope
}

func (r Register) Username() string {
	return r.username
}

func (r Register) Password() string {
	return r.password
}

func (r Register) FullName() string {
	return r.fullName
}
