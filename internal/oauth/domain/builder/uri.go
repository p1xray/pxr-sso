package builder

import (
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"net/url"
	"strings"
)

type URI struct {
	loginRedirectURI        string
	consentRedirectURI      string
	defaultErrorRedirectURI string
}

func NewURI() *URI {
	return &URI{
		// TODO: set this from config
		loginRedirectURI:        "http://localhost:3000/login",
		consentRedirectURI:      "http://localhost:3000/consent",
		defaultErrorRedirectURI: "http://localhost:3000/sigin/error",
	}
}

func (u *URI) BuildLoginRedirectURI(flow dto.Flow) string {
	queryValues := map[string]string{
		oauth.RequestParameterNameFlowID:       flow.ID().String(),
		oauth.RequestParameterNameResponseType: flow.ResponseType(),
		oauth.RequestParameterNameClientID:     flow.ClientID(),
		oauth.RequestParameterNameRedirectURI:  flow.RedirectURI(),
		oauth.RequestParameterNameState:        flow.State(),
		oauth.RequestParameterNameScope:        strings.Join(flow.Scope(), " "),
	}

	redirectURI := u.buildRedirectURI(u.loginRedirectURI, queryValues)
	return redirectURI
}

func (u *URI) BuildConsentRedirectURI(flow dto.Flow) string {
	queryValues := map[string]string{
		oauth.RequestParameterNameFlowID:       flow.ID().String(),
		oauth.RequestParameterNameResponseType: flow.ResponseType(),
		oauth.RequestParameterNameClientID:     flow.ClientID(),
		oauth.RequestParameterNameRedirectURI:  flow.RedirectURI(),
		oauth.RequestParameterNameState:        flow.State(),
		oauth.RequestParameterNameScope:        strings.Join(flow.Scope(), " "),
	}

	redirectURI := u.buildRedirectURI(u.consentRedirectURI, queryValues)
	return redirectURI
}

func (u *URI) BuildCallbackRedirectURI(flow dto.Flow) string {
	queryValues := map[string]string{
		oauth.RequestParameterNameAuthorizationCode: flow.AuthorizationCode(),
		oauth.RequestParameterNameState:             flow.State(),
	}

	redirectURI := u.buildRedirectURI(flow.RedirectURI(), queryValues)
	return redirectURI
}

func (u *URI) BuildErrorRedirectURI(rawURL string, oauthErr *domain.OAuthError) string {
	// build error redirect URI on redirect URI from request parameters
	errorRedirectURI := u.buildErrorRedirectURIByOAuthError(rawURL, oauthErr)

	if errorRedirectURI == "" {
		// otherwise build error redirect URI on authorize service default error page
		errorRedirectURI = u.buildErrorRedirectURIByOAuthError(u.defaultErrorRedirectURI, oauthErr)
	}

	return errorRedirectURI
}

func (u *URI) buildErrorRedirectURIByOAuthError(rawURL string, oauthErr *domain.OAuthError) string {
	if rawURL == "" {
		return ""
	}

	if oauthErr == nil {
		// build redirect URI with internal server error
		return u.buildErrorRedirectURI(rawURL, domain.ErrorCodeServerError, domain.ErrorDescriptionInternalServerError, "")
	}

	return u.buildErrorRedirectURI(rawURL, oauthErr.Code, oauthErr.Description, oauthErr.URI)
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
