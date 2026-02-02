package domain

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
	ErrOAuthClientNotRegistered               = errors.New("the client with the provided client_id is not registered")
)

// OAuthError wraps OAuth errors with additional context.
// It provides structured error information that used for
// returning appropriate error responses.
type OAuthError struct {
	// Code is a machine-readable error code (e.g., "invalid_request", "unauthorized_client")
	Code string

	// Message is a human-readable error message
	Description string

	// URI identifying a human-readable web page with information about the error
	URI string

	// details contains the underlying error
	details error
}

func newOAuthErrorError(code string, description string, uri string, err error) *OAuthError {
	return &OAuthError{Code: code, Description: description, URI: uri, details: err}
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

func (e *OAuthError) Unwrap() error {
	return e.details
}

func InvalidRequestOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeInvalidRequest, err.Error(), "", err)
}

func UnauthorizedClientOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeUnauthorizedClient, err.Error(), "", err)
}

func UnsupportedResponseTypeOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeUnsupportedResponseType, err.Error(), "", err)
}

func AccessDeniedOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeAccessDenied, err.Error(), "", err)
}

func InvalidScopeOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeInvalidScope, err.Error(), "", err)
}

func ServerErrorOAuthError(err error) *OAuthError {
	return newOAuthErrorError(ErrorCodeServerError, ErrorDescriptionInternalServerError, "", err)
}

// Common error codes
const (
	ErrorCodeInvalidRequest          = "invalid_request"
	ErrorCodeUnauthorizedClient      = "unauthorized_client"
	ErrorCodeUnsupportedResponseType = "unsupported_response_type"
	ErrorCodeAccessDenied            = "access_denied"
	ErrorCodeInvalidScope            = "invalid_scope"
	ErrorCodeServerError             = "server_error"
)

// Common error descriprions
const (
	ErrorDescriptionInternalServerError = "internal server error"
)
