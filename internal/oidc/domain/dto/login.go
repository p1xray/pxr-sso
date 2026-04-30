package dto

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
)

type LoginRequest struct {
	requestURI string
	username   string
	password   string
}

func NewLoginRequest(
	requestURI,
	username,
	password string,
) LoginRequest {
	return LoginRequest{
		requestURI: requestURI,
		username:   username,
		password:   password,
	}
}

func (l *LoginRequest) RequestURI() string {
	return l.requestURI
}

func (l *LoginRequest) Username() string {
	return l.username
}

func (l *LoginRequest) Password() string {
	return l.password
}

func (l *LoginRequest) Validate() error {
	if l.RequestURI() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameRequestURI)
	}

	if l.Username() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameUsername)
	}

	if l.Password() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNamePassword)
	}

	return nil
}

type LoginResponse struct {
	redirectURI string
	session     SessionCookie
}

func NewLoginResponse(redirectURI string, session SessionCookie) LoginResponse {
	return LoginResponse{
		redirectURI: redirectURI,
		session:     session,
	}
}

func (l *LoginResponse) RedirectURI() string {
	return l.redirectURI
}

func (l *LoginResponse) Session() SessionCookie {
	return l.session
}
