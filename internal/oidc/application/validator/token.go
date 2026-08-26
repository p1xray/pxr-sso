package validator

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/hasher/sha256"
	"slices"
)

type tokenRequestValidator struct {
	tokenRequest     dto.TokenRequest
	authorizeRequest dto.ValidatedAuthorizeRequest
	client           dto.Client
}

func NewTokenRequestValidator(
	tokenRequest dto.TokenRequest,
	authorizeRequest dto.ValidatedAuthorizeRequest,
	client dto.Client,
) *tokenRequestValidator {
	return &tokenRequestValidator{
		tokenRequest:     tokenRequest,
		authorizeRequest: authorizeRequest,
		client:           client,
	}
}

func (v *tokenRequestValidator) Validate() error {
	if err := v.validateGrantType(); err != nil {
		return err
	}

	if err := v.validateClientID(); err != nil {
		return err
	}

	if err := v.validateAuthorizationCode(); err != nil {
		return err
	}

	if err := v.validateRedirectURI(); err != nil {
		return err
	}

	if err := v.validateCodeVerifier(); err != nil {
		return err
	}

	return nil
}

func (v *tokenRequestValidator) validateGrantType() *oidc.Error {
	if err := v.validateGrantTypeRequired(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	if err := v.validateGrantTypeValue(); err != nil {
		return oidc.UnsupportedGrantTypeError(err)
	}

	return nil
}

func (v *tokenRequestValidator) validateGrantTypeRequired() error {
	if v.tokenRequest.GrantType() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameGrantType)
	}

	return nil
}

func (v *tokenRequestValidator) validateGrantTypeValue() error {
	if v.tokenRequest.GrantType() != RequestParameterAllowValueGrantTypeAuthorizationCode {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameGrantType)
	}

	return nil
}

func (v *tokenRequestValidator) validateClientID() *oidc.Error {
	if err := v.validateClientIDRequired(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	if err := v.validateClientIDExistClient(); err != nil {
		return oidc.UnauthorizedClientError(err)
	}

	return nil
}

func (v *tokenRequestValidator) validateClientIDRequired() error {
	if v.tokenRequest.ClientID() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameClientID)
	}

	return nil
}

func (v *tokenRequestValidator) validateClientIDExistClient() error {
	if v.tokenRequest.ClientID() != v.client.Code() {
		return oidc.ErrOAuthClientNotRegistered
	}

	return nil
}

func (v *tokenRequestValidator) validateAuthorizationCode() *oidc.Error {
	if err := v.validateAuthorizationCodeRequired(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	return nil
}

func (v *tokenRequestValidator) validateAuthorizationCodeRequired() error {
	if v.tokenRequest.AuthorizationCode() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameAuthorizationCode)
	}

	return nil
}

func (v *tokenRequestValidator) validateRedirectURI() *oidc.Error {
	if err := v.validateRedirectURIRequired(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		return oidc.UnauthorizedClientError(err)
	}

	if err := v.validateRedirectURIValue(); err != nil {
		return oidc.UnauthorizedClientError(err)
	}

	return nil
}

func (v *tokenRequestValidator) validateRedirectURIRequired() error {
	if v.tokenRequest.RedirectURI() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *tokenRequestValidator) validateRedirectURIValue() error {
	if v.tokenRequest.RedirectURI() != v.authorizeRequest.RedirectURI() {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *tokenRequestValidator) validateRedirectURIRegisteredForClient() error {
	if v.redirectURIRegisteredForClient() == false {
		return oidc.ErrOAuthRedirectURINotRegisteredForClient
	}

	return nil
}

func (v *tokenRequestValidator) redirectURIRegisteredForClient() bool {
	return slices.Contains(v.client.RedirectURI(), v.tokenRequest.RedirectURI())
}

func (v *tokenRequestValidator) validateCodeVerifier() *oidc.Error {
	if err := v.validateCodeVerifierRequired(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	if err := v.validateCodeVerifierValue(); err != nil {
		return oidc.InvalidRequestError(err)
	}

	return nil
}

func (v *tokenRequestValidator) validateCodeVerifierRequired() error {
	if v.tokenRequest.CodeVerifier() == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameCodeVerifier)
	}

	return nil
}

func (v *tokenRequestValidator) validateCodeVerifierValue() error {
	hashCodeVerifier := sha256.Base64URLHash(v.tokenRequest.CodeVerifier())
	codeChallenge := v.authorizeRequest.CodeChallenge()
	if hashCodeVerifier != codeChallenge {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameCodeVerifier)
	}

	return nil
}
