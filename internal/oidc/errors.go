package oidc

import (
	"errors"
	"fmt"
)

// errors for OAuth request parameters validation.
var (
	ErrOAuthMissingRequiredParameter          = errors.New("the request is missing the required parameter")
	ErrOAuthParameterIncludedMoreThanOnce     = errors.New("the parameter included in the request more than once")
	ErrOAuthParameterInvalidValue             = errors.New("the parameter have an invalid value")
	ErrOAuthRedirectURINotRegisteredForClient = errors.New("the redirect_uri used in the request is not registered for the client being used")
	ErrOAuthAudienceNotRegisteredForClient    = errors.New("the audience used in the request is not registered for the client being used")
	ErrOAuthClientNotRegistered               = errors.New("the client with the provided client_id is not registered")
	ErrOAuthFlowNotExists                     = errors.New("the flow with the provided flow_id is not exists")
	ErrOAuthAuthorizationNotExists            = errors.New("the authorization with the provided code is not exists")
	ErrOAuthUsernameRequired                  = errors.New("username is required")
	ErrOAuthPasswordRequired                  = errors.New("password is required")
	ErrOAuthFullNameRequired                  = errors.New("full name is required")
	ErrOAuthInvalidUserCredentials            = errors.New("invalid username or password")
	ErrOAuthUserAlreadyExists                 = errors.New("user with same username is already exists")
)

// errors for OIDC protocol.
var (
	ErrOIDCLoginRequired            = errors.New("the authorization server requires end-user authentication")
	ErrOIDCAccountSelectionRequired = errors.New("the end-user is to select a session at the authorization server")
	ErrOIDCConsentRequired          = errors.New("the authorization server requires end-user consent")
)

// Error codes.
const (
	ErrorCodeServerError                = "server_error"
	ErrorDescriptionInternalServerError = "internal server error"

	errorCodeInvalidRequest           = "invalid_request"
	errorCodeInvalidClient            = "invalid_client"
	errorCodeUnauthorizedClient       = "unauthorized_client"
	errorCodeUnsupportedResponseType  = "unsupported_response_type"
	errorCodeAccessDenied             = "access_denied"
	errorCodeInvalidScope             = "invalid_scope"
	errorCodeInvalidUserCredentials   = "invalid_user_credentials"
	errorCodeInvalidGrant             = "invalid_grant"
	errorCodeUnsupportedGrantType     = "unsupported_grant_type"
	errorCodeLoginRequired            = "login_required"
	errorCodeAccountSelectionRequired = "account_selection_required"
	errorCodeConsentRequired          = "consent_required"
)

// Error wraps OAuth errors with additional context.
// It provides structured error information that used for
// returning appropriate error responses.
type Error struct {
	// Code is a machine-readable error code (e.g., "invalid_request", "unauthorized_client").
	Code string

	// Message is a human-readable error message.
	Description string

	// URI identifying a human-readable web page with information about the error.
	URI string

	// details contains the underlying error.
	details error
}

func (e *Error) Error() string {
	description := ""
	if e.Description != "" {
		description = ": " + e.Description
	}

	errMessage := ": " + e.details.Error()

	return fmt.Sprintf("%s: %s%s", e.Code, description, errMessage)
}

func (e *Error) Unwrap() error {
	return e.details
}

func newError(code string, description string, uri string, err error) *Error {
	return &Error{Code: code, Description: description, URI: uri, details: err}
}

func InvalidRequestError(err error) *Error {
	return newError(errorCodeInvalidRequest, err.Error(), "", err)
}

func InvalidClientError(err error) *Error {
	return newError(errorCodeInvalidClient, err.Error(), "", err)
}

func UnauthorizedClientError(err error) *Error {
	return newError(errorCodeUnauthorizedClient, err.Error(), "", err)
}

func UnsupportedResponseTypeError(err error) *Error {
	return newError(errorCodeUnsupportedResponseType, err.Error(), "", err)
}

func AccessDeniedError(err error) *Error {
	return newError(errorCodeAccessDenied, err.Error(), "", err)
}

func InvalidScopeError(err error) *Error {
	return newError(errorCodeInvalidScope, err.Error(), "", err)
}

func ServerError(err error) *Error {
	return newError(ErrorCodeServerError, ErrorDescriptionInternalServerError, "", err)
}

func InvalidUserCredentialsError() *Error {
	return newError(errorCodeInvalidUserCredentials, ErrOAuthInvalidUserCredentials.Error(), "", ErrOAuthInvalidUserCredentials)
}

func InvalidGrantError(err error) *Error {
	return newError(errorCodeInvalidGrant, err.Error(), "", err)
}

func UnsupportedGrantTypeError(err error) *Error {
	return newError(errorCodeUnsupportedGrantType, err.Error(), "", err)
}

func LoginRequiredError() *Error {
	err := ErrOIDCLoginRequired
	return newError(errorCodeLoginRequired, err.Error(), "", err)
}

func AccountSelectionRequiredError() *Error {
	err := ErrOIDCAccountSelectionRequired
	return newError(errorCodeAccountSelectionRequired, err.Error(), "", err)
}

func ConsentRequiredError() *Error {
	err := ErrOIDCConsentRequired
	return newError(errorCodeConsentRequired, err.Error(), "", err)
}

func (e *Error) IsInvalidUserCredentials() bool {
	return e.Code == errorCodeInvalidUserCredentials
}
