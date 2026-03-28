package token

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator"
	"github.com/p1xray/pxr-sso/pkg/hasher/sha256"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"slices"
)

type Validator struct {
	params        dto.ExchangeToken
	client        nullable.Nullable[dto.Client]
	authorization nullable.Nullable[dto.Authorization]
}

func NewValidator(
	params dto.ExchangeToken,
	client nullable.Nullable[dto.Client],
	authorization nullable.Nullable[dto.Authorization],
) *Validator {
	return &Validator{
		params:        params,
		client:        client,
		authorization: authorization,
	}
}

func (v *Validator) Validate() error {
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

func (v *Validator) validateGrantType() *oauth.Error {
	if err := v.validateGrantTypeRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateGrantTypeValue(); err != nil {
		return oauth.InvalidUnsupportedGrantTypeError(err)
	}

	return nil
}

func (v *Validator) validateGrantTypeRequired() error {
	if v.params.GrantType() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameGrantType)
	}

	return nil
}

func (v *Validator) validateGrantTypeValue() error {
	if v.params.GrantType() != validator.RequestParameterAllowValueGrantTypeAuthorizationCode {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameGrantType)
	}

	return nil
}

func (v *Validator) validateClientID() *oauth.Error {
	if err := v.validateClientIDRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateClientIDExistClient(); err != nil {
		return oauth.UnauthorizedClientError(err)
	}

	return nil
}

func (v *Validator) validateClientIDRequired() error {
	if v.params.ClientID() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDExistClient() error {
	if v.client.IsNone() {
		return oauth.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()

	if v.params.ClientID() != client.Code() {
		return oauth.ErrOAuthClientNotRegistered
	}

	return nil
}

func (v *Validator) validateAuthorizationCode() *oauth.Error {
	if err := v.validateAuthorizationCodeRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateAuthorizationCodeValue(); err != nil {
		return oauth.InvalidGrantError(err)
	}

	return nil
}

func (v *Validator) validateAuthorizationCodeRequired() error {
	if v.params.AuthorizationCode() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameAuthorizationCode)
	}

	return nil
}

func (v *Validator) validateAuthorizationCodeValue() error {
	if v.authorization.IsNone() {
		return oauth.ErrOAuthAuthorizationNotExists
	}

	authorization := v.authorization.Unwrap()
	if v.params.AuthorizationCode() != authorization.Code() {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameAuthorizationCode)
	}

	return nil
}

func (v *Validator) validateRedirectURI() *oauth.Error {
	if err := v.validateRedirectURIRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		return oauth.UnauthorizedClientError(err)
	}

	if err := v.validateRedirectURIValue(); err != nil {
		return oauth.UnauthorizedClientError(err)
	}

	return nil
}

func (v *Validator) validateRedirectURIRequired() error {
	if v.params.RedirectURI() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIValue() error {
	if v.authorization.IsNone() {
		return oauth.ErrOAuthAuthorizationNotExists
	}

	authorization := v.authorization.Unwrap()
	if v.params.RedirectURI() != authorization.RedirectURI() {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIRegisteredForClient() error {
	if v.redirectURIRegisteredForClient() == false {
		return oauth.ErrOAuthRedirectURINotRegisteredForClient
	}

	return nil
}

func (v *Validator) redirectURIRegisteredForClient() bool {
	if v.client.IsNone() {
		return false
	}

	client := v.client.Unwrap()
	return slices.Contains(client.RedirectURI(), v.params.RedirectURI())
}

func (v *Validator) validateCodeVerifier() *oauth.Error {
	if err := v.validateCodeVerifierRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateCodeVerifierValue(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	return nil
}

func (v *Validator) validateCodeVerifierRequired() error {
	if v.params.CodeVerifier() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameCodeVerifier)
	}

	return nil
}

func (v *Validator) validateCodeVerifierValue() error {
	if v.authorization.IsNone() {
		return oauth.ErrOAuthAuthorizationNotExists
	}

	authorization := v.authorization.Unwrap()
	hashCodeVerifier := sha256.Base64URLHash(v.params.CodeVerifier())
	codeChallenge := authorization.CodeChallenge()
	if hashCodeVerifier != codeChallenge {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameCodeVerifier)
	}

	return nil
}
