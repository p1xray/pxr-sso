package dto

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
)

type RegisterRequest struct {
	requestURI string
	username   string
	password   string
	fullName   string
}

func NewRegisterRequest(
	requestURI,
	username,
	password,
	fullName string,
) RegisterRequest {
	return RegisterRequest{
		requestURI: requestURI,
		username:   username,
		password:   password,
		fullName:   fullName,
	}
}

func (r *RegisterRequest) RequestURI() string {
	return r.requestURI
}

func (r *RegisterRequest) Username() string {
	return r.username
}

func (r *RegisterRequest) Password() string {
	return r.password
}

func (r *RegisterRequest) FullName() string {
	return r.fullName
}

func (r *RegisterRequest) Validate() error {
	if r.RequestURI() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameRequestURI)
	}

	if r.Username() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameUsername)
	}

	if r.Password() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNamePassword)
	}

	if r.FullName() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameFullName)
	}

	return nil
}

type RegisterResponse struct {
	redirectURI string
	session     SessionCookie
}

func NewRegisterResponse(redirectURI string, session SessionCookie) RegisterResponse {
	return RegisterResponse{
		redirectURI: redirectURI,
		session:     session,
	}
}

func (r *RegisterResponse) RedirectURI() string {
	return r.redirectURI
}

func (r *RegisterResponse) Session() SessionCookie {
	return r.session
}
