package dto

type Consent struct {
	flowID       string
	responseType string
	clientID     string
	redirectURI  string
	state        string
	scope        []string
}

func NewConsent(
	flowID,
	responseType,
	clientID,
	redirectURI,
	state string,
	scope []string,
) Consent {
	return Consent{
		flowID:       flowID,
		responseType: responseType,
		clientID:     clientID,
		redirectURI:  redirectURI,
		state:        state,
		scope:        scope,
	}
}

func (c *Consent) FlowID() string {
	return c.flowID
}

func (c *Consent) ResponseType() string {
	return c.responseType
}

func (c *Consent) ClientID() string {
	return c.clientID
}

func (c *Consent) RedirectURI() string {
	return c.redirectURI
}

func (c *Consent) State() string {
	return c.state
}

func (c *Consent) Scope() []string {
	return c.scope
}
