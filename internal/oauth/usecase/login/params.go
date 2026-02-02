package login

type Params struct {
	FlowID       string
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
	Scope        string
	Username     string
	Password     string
}
