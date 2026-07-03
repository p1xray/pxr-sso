package dto

type ValidatedAuthorizeRequest struct {
	responseType        string
	prompt              string
	clientID            string
	redirectURI         string
	codeChallenge       string
	codeChallengeMethod string
	state               string
	audience            string
	scopes              ValidatedScopesRequest
}

func NewValidatedAuthorizeRequest(
	responseType,
	prompt,
	clientID,
	redirectURI,
	codeChallenge,
	codeChallengeMethod,
	state,
	audience string,
	scopes ValidatedScopesRequest,
) ValidatedAuthorizeRequest {
	return ValidatedAuthorizeRequest{
		responseType:        responseType,
		prompt:              prompt,
		clientID:            clientID,
		redirectURI:         redirectURI,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		state:               state,
		audience:            audience,
		scopes:              scopes,
	}
}

func (v *ValidatedAuthorizeRequest) ResponseType() string {
	return v.responseType
}

func (v *ValidatedAuthorizeRequest) Prompt() string {
	return v.prompt
}

func (v *ValidatedAuthorizeRequest) ClientID() string {
	return v.clientID
}

func (v *ValidatedAuthorizeRequest) RedirectURI() string {
	return v.redirectURI
}

func (v *ValidatedAuthorizeRequest) CodeChallenge() string {
	return v.codeChallenge
}

func (v *ValidatedAuthorizeRequest) CodeChallengeMethod() string {
	return v.codeChallengeMethod
}

func (v *ValidatedAuthorizeRequest) State() string {
	return v.state
}

func (v *ValidatedAuthorizeRequest) Audience() string {
	return v.audience
}

func (v *ValidatedAuthorizeRequest) Scopes() ValidatedScopesRequest {
	return v.scopes
}

func (v *ValidatedAuthorizeRequest) GrantedScopes() []string {
	return v.scopes.Granted()
}

func (v *ValidatedAuthorizeRequest) PendingScopes() []string {
	return v.scopes.Pending()
}

func (v *ValidatedAuthorizeRequest) AllScopes() []string {
	return v.scopes.All()
}

type ValidatedScopesRequest struct {
	granted []string
	pending []string
}

type ValidatedScopesRequestOption func(*ValidatedScopesRequest)

func NewValidatedScopesRequest(opts ...ValidatedScopesRequestOption) ValidatedScopesRequest {
	validatedScopesRequest := ValidatedScopesRequest{
		granted: make([]string, 0),
		pending: make([]string, 0),
	}

	for _, opt := range opts {
		opt(&validatedScopesRequest)
	}

	return validatedScopesRequest
}

func (v *ValidatedScopesRequest) Granted() []string {
	return v.granted
}

func (v *ValidatedScopesRequest) Pending() []string {
	return v.pending
}

func (v *ValidatedScopesRequest) All() []string {
	allScopes := make([]string, 0, len(v.granted)+len(v.pending))
	allScopes = append(allScopes, v.granted...)
	allScopes = append(allScopes, v.pending...)

	return allScopes
}

func WithGranted(granted []string) ValidatedScopesRequestOption {
	return func(v *ValidatedScopesRequest) {
		v.granted = granted
	}
}

func WithPending(pending []string) ValidatedScopesRequestOption {
	return func(v *ValidatedScopesRequest) {
		v.pending = pending
	}
}
