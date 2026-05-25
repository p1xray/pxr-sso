package dto

type AuthorizedGrant struct {
	session AuthorizedSession
	request ValidatedAuthorizeRequest
}

func NewAuthorizedGrant(session AuthorizedSession, request ValidatedAuthorizeRequest) AuthorizedGrant {
	return AuthorizedGrant{
		session: session,
		request: request,
	}
}

func (a *AuthorizedGrant) Session() AuthorizedSession {
	return a.session
}

func (a *AuthorizedGrant) Request() ValidatedAuthorizeRequest {
	return a.request
}

func (a *AuthorizedGrant) RequestClientID() string {
	return a.request.ClientID()
}

func (a *AuthorizedGrant) SessionUserID() int64 {
	return a.session.Subject()
}

func (a *AuthorizedGrant) RequestGrantedScopes() []string {
	return a.request.GrantedScopes()
}

func (a *AuthorizedGrant) RequestAudience() string {
	return a.request.Audience()
}
