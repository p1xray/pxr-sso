package builder

import (
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"net/url"
)

type URI struct {
	defaultErrorRedirectURI string
	loginRedirectURI        string
}

func NewURI() *URI {
	return &URI{
		// TODO: set this from config
		loginRedirectURI:        "http://localhost:3000/login",
		defaultErrorRedirectURI: "http://localhost:3000/sigin/error",
	}
}

func (u *URI) BuildLoginRedirectURI(flow dto.Flow) string {
	queryValues := map[string]string{
		oauth.RequestParameterNameFlowID:      flow.ID().String(),
		oauth.RequestParameterNameClientID:    flow.ClientID(),
		oauth.RequestParameterNameRedirectURI: flow.RedirectURI(),
		oauth.RequestParameterNameState:       flow.State(),
	}

	redirectURI := u.buildRedirectURI(u.loginRedirectURI, queryValues)

	return redirectURI
}

func (u *URI) BuildErrorRedirectURI(paramRedirectURI string, oauthErr *oauth.OAuthError) string {
	// build error redirect URI on redirect URI from request parameters
	errorRedirectURI := u.buildErrorRedirectURIByOAuthError(paramRedirectURI, oauthErr)

	if errorRedirectURI == "" {
		// otherwise build error redirect URI on authorize service default error page
		errorRedirectURI = u.buildErrorRedirectURIByOAuthError(u.defaultErrorRedirectURI, oauthErr)
	}

	return errorRedirectURI
}

func (u *URI) buildErrorRedirectURIByOAuthError(redirectURI string, oauthErr *oauth.OAuthError) string {
	if redirectURI == "" {
		return ""
	}

	if oauthErr == nil {
		// build redirect URI with internal server error
		return u.buildErrorRedirectURI(redirectURI, oauth.ErrorCodeServerError, oauth.ErrorDescriptionInternalServerError, "")
	}

	return u.buildErrorRedirectURI(redirectURI, oauthErr.Code, oauthErr.Description, oauthErr.URI)
}

func (u *URI) buildErrorRedirectURI(rawURL, code, description, errURI string) string {
	queryValues := map[string]string{
		oauth.RequestParameterNameError:            code,
		oauth.RequestParameterNameErrorDescription: description,
		oauth.RequestParameterNameErrorURI:         errURI,
	}

	redirectURI := u.buildRedirectURI(rawURL, queryValues)

	return redirectURI
}

func (u *URI) buildRedirectURI(rawURL string, queryValues map[string]string) string {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	query := uri.Query()
	for k, v := range queryValues {
		query.Set(k, v)
	}

	uri.RawQuery = query.Encode()
	return uri.String()
}
