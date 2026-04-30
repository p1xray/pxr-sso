package dto

type Authorization struct {
	code        string
	sessionID   string
	clientID    string
	userID      int64
	redirectURI string
	state       string
	audience    string
	scope       []string
}

func NewAuthorization(
	code,
	sessionID,
	clientID,
	redirectURI,
	state,
	audience string,
	scope []string,
	userID int64,
) Authorization {
	return Authorization{
		code:        code,
		sessionID:   sessionID,
		clientID:    clientID,
		userID:      userID,
		redirectURI: redirectURI,
		state:       state,
		audience:    audience,
		scope:       scope,
	}
}

func (a *Authorization) Code() string {
	return a.code
}

func (a *Authorization) SessionID() string {
	return a.sessionID
}

func (a *Authorization) ClientID() string {
	return a.clientID
}

func (a *Authorization) UserID() int64 {
	return a.userID
}

func (a *Authorization) RedirectURI() string {
	return a.redirectURI
}

func (a *Authorization) State() string {
	return a.state
}

func (a *Authorization) Audience() string {
	return a.audience
}

func (a *Authorization) Scope() []string {
	return a.scope
}
