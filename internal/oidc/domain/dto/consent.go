package dto

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
)

type ConsentRequest struct {
	requestURI     string
	scopes         []string
	sessionCookies []SessionCookie
}

func NewConsentRequest(
	requestURI string,
	scopes []string,
	sessionCookies []SessionCookie,
) ConsentRequest {
	return ConsentRequest{
		requestURI:     requestURI,
		scopes:         scopes,
		sessionCookies: sessionCookies,
	}
}

func (c *ConsentRequest) RequestURI() string {
	return c.requestURI
}

func (c *ConsentRequest) Scopes() []string {
	return c.scopes
}

func (c *ConsentRequest) SessionCookies() []SessionCookie {
	return c.sessionCookies
}

func (c *ConsentRequest) Validate() error {
	if c.RequestURI() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameRequestURI)
	}

	if len(c.SessionCookies()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameSession)
	}

	if len(c.SessionCookies()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameSession)
	}

	return nil
}

func (c *ConsentRequest) SingleSessionCookie() (SessionCookie, error) {
	const op = "get single session cookie"

	if len(c.SessionCookies()) == 0 {
		return SessionCookie{}, fmt.Errorf("%s: no sessions found in consent request", op)
	}

	if len(c.SessionCookies()) > 1 {
		return SessionCookie{}, fmt.Errorf("%s: multiple sessions found in consent request", op)
	}

	return c.SessionCookies()[0], nil
}

type ConsentCardRequest struct {
	requestURI string
}

func NewConsentCardRequest(requestURI string) ConsentCardRequest {
	return ConsentCardRequest{
		requestURI: requestURI,
	}
}

func (c *ConsentCardRequest) RequestURI() string {
	return c.requestURI
}

type ConsentCardResponse struct {
	scopes []ConsentScope
}

func NewConsentCardResponse(scopes []ConsentScope) ConsentCardResponse {
	return ConsentCardResponse{
		scopes: scopes,
	}
}

func (c *ConsentCardResponse) Scopes() []ConsentScope {
	return c.scopes
}

type ConsentScope struct {
	code        string
	name        string
	description string
	isGranted   bool
}

func NewConsentScope(code, name, description string, isGranted bool) ConsentScope {
	return ConsentScope{
		code:        code,
		name:        name,
		description: description,
		isGranted:   isGranted,
	}
}

func (c *ConsentScope) Code() string {
	return c.code
}

func (c *ConsentScope) Name() string {
	return c.name
}

func (c *ConsentScope) Description() string {
	return c.description
}

func (c *ConsentScope) IsGranted() bool {
	return c.isGranted
}
