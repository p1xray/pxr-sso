package builder

import (
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator"
	"net/url"
	"strings"
)

type URI struct {
	loginRedirectURI        string
	consentRedirectURI      string
	defaultErrorRedirectURI string
}

func NewURI(cfg URIBuilderConfig) *URI {
	return &URI{
		loginRedirectURI:        cfg.Login,
		consentRedirectURI:      cfg.Consent,
		defaultErrorRedirectURI: cfg.Error,
	}
}

func (u *URI) BuildLoginRedirectURI(flow dto.Flow) string {
	queryValues := newLoginQueryParameters(
		flow.ID().String(),
		flow.ResponseType(),
		flow.ClientID(),
		flow.RedirectURI(),
		flow.State(),
		strings.Join(flow.Scope(), " "),
	)

	redirectURI := u.buildRedirectURI(u.loginRedirectURI, queryValues)
	return redirectURI
}

func (u *URI) BuildConsentRedirectURI(flow dto.Flow) string {
	queryValues := newConsentQueryParameters(
		flow.ID().String(),
		flow.ResponseType(),
		flow.ClientID(),
		flow.RedirectURI(),
		flow.State(),
		strings.Join(flow.Scope(), " "),
	)

	redirectURI := u.buildRedirectURI(u.consentRedirectURI, queryValues)
	return redirectURI
}

func (u *URI) BuildCallbackRedirectURI(flow dto.Flow) string {
	queryValues := newCallbackQueryParameters(flow.AuthorizationCode(), flow.State())

	redirectURI := u.buildRedirectURI(flow.RedirectURI(), queryValues)
	return redirectURI
}

func (u *URI) BuildErrorRedirectURI(rawURL string, err error) string {
	baseRedirectURI := u.fetchErrorRedirectURI(rawURL)
	parameters := u.fetchErrorQueryParameters(err)

	errorRedirectURI := u.buildRedirectURI(baseRedirectURI, parameters)
	return errorRedirectURI
}

func (u *URI) fetchErrorRedirectURI(rawURL string) string {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return u.defaultErrorRedirectURI
	}

	return rawURL
}

func (u *URI) fetchErrorQueryParameters(err error) queryParameters {
	var validationErr *validator.Error
	if errors.As(err, &validationErr) {
		return newErrorQueryParameters(validationErr.Code, validationErr.Description, validationErr.URI)
	}

	return newErrorQueryParameters(domain.ErrorCodeServerError, domain.ErrorDescriptionInternalServerError, "")
}

func (u *URI) buildRedirectURI(rawURL string, parameters queryParameters) string {
	uri, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ""
	}

	query := uri.Query()
	for k, v := range parameters {
		query.Set(k, v)
	}

	uri.RawQuery = query.Encode()
	return uri.String()
}
