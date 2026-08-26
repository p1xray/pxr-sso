package builder

import (
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"net/url"
)

type URIBuilder interface {
	BuildLoginRedirectURI(requestURI string) (string, error)
	BuildConsentRedirectURI(requestURI string) (string, error)
	BuildSelectAccountRedirectURI(requestURI string) (string, error)
	BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error)
	BuildErrorRedirectURI(rawURL string, err error) (string, error)
}

type uri struct {
	loginRedirectURI         string
	consentRedirectURI       string
	selectAccountRedirectURI string
	defaultErrorRedirectURI  string
}

func NewURIBuilder(cfg URIBuilderConfig) *uri {
	return &uri{
		loginRedirectURI:         cfg.Login,
		consentRedirectURI:       cfg.Consent,
		selectAccountRedirectURI: cfg.SelectAccount,
		defaultErrorRedirectURI:  cfg.Error,
	}
}

func (u *uri) BuildLoginRedirectURI(requestURI string) (string, error) {
	queryValues := newRequestURIQueryParameters(requestURI)

	return u.buildRedirectURI(u.loginRedirectURI, queryValues)
}

func (u *uri) BuildConsentRedirectURI(requestURI string) (string, error) {
	queryValues := newRequestURIQueryParameters(requestURI)

	return u.buildRedirectURI(u.consentRedirectURI, queryValues)
}

func (u *uri) BuildSelectAccountRedirectURI(requestURI string) (string, error) {
	queryValues := newRequestURIQueryParameters(requestURI)

	return u.buildRedirectURI(u.selectAccountRedirectURI, queryValues)
}

func (u *uri) BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error) {
	queryValues := newCallbackQueryParameters(authorizationCode, state)

	return u.buildRedirectURI(rawURL, queryValues)
}

func (u *uri) BuildErrorRedirectURI(rawURL string, err error) (string, error) {
	baseRedirectURI := u.fetchErrorRedirectURI(rawURL)
	parameters := u.fetchErrorQueryParameters(err)

	return u.buildRedirectURI(baseRedirectURI, parameters)
}

func (u *uri) fetchErrorRedirectURI(rawURL string) string {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return u.defaultErrorRedirectURI
	}

	return rawURL
}

func (u *uri) fetchErrorQueryParameters(err error) queryParameters {
	var validationErr *oidc.Error
	if errors.As(err, &validationErr) {
		return newErrorQueryParameters(validationErr.Code, validationErr.Description, validationErr.URI)
	}

	return newErrorQueryParameters(oidc.ErrorCodeServerError, oidc.ErrorDescriptionInternalServerError, "")
}

func (u *uri) buildRedirectURI(rawURL string, parameters queryParameters) (string, error) {
	uri, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse raw url: %w", err)
	}

	query := uri.Query()
	for k, v := range parameters {
		query.Set(k, v)
	}

	uri.RawQuery = query.Encode()
	return uri.String(), nil
}
