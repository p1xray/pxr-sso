package token

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/hasher/sha256"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"slices"
)

type Validator struct {
	params dto.ExchangeToken
	client nullable.Nullable[dto.Client]
	flow   nullable.Nullable[dto.Flow]
}

func NewValidator(
	params dto.ExchangeToken,
	client nullable.Nullable[dto.Client],
	flow nullable.Nullable[dto.Flow],
) *Validator {
	return &Validator{
		params: params,
		client: client,
		flow:   flow,
	}
}

func (v *Validator) Validate() error {
	if err := v.validateFlowID(); err != nil {
		return err
	}

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

	if err := v.validateAudience(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidatedScope() []string {
	if v.flow.IsNone() {
		return []string{}
	}

	flow := v.flow.Unwrap()
	scope := extslices.Intersect(flow.Scope(), v.params.Scope())

	return scope
}

func (v *Validator) validateFlowID() *oauth.Error {
	if err := v.validateFlowIDRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateFlowIDValue(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateFlowIDExistFlow(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	return nil
}

func (v *Validator) validateFlowIDRequired() error {
	if v.params.FlowID() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDValue() error {
	if err := uuid.Validate(v.params.FlowID()); err != nil {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDExistFlow() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()

	paramFlowID, err := uuid.Parse(v.params.FlowID())
	if err != nil {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	if paramFlowID != flow.ID() {
		return oauth.ErrOAuthFlowNotExists
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

	if err := v.validateAuthorizationCodeEqualsFlowAuthorizationCode(); err != nil {
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

func (v *Validator) validateAuthorizationCodeEqualsFlowAuthorizationCode() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.AuthorizationCode() != flow.AuthorizationCode() {
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

	if err := v.validateRedirectURIEqualsFlowRedirectURI(); err != nil {
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

func (v *Validator) validateRedirectURIEqualsFlowRedirectURI() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.RedirectURI() != flow.RedirectURI() {
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

	if err := v.validateCodeVerifierHashEqualsFlowCodeChallenge(); err != nil {
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

func (v *Validator) validateCodeVerifierHashEqualsFlowCodeChallenge() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	hashCodeVerifier := sha256.Base64URLHash(v.params.CodeVerifier())
	codeChallenge := flow.CodeChallenge()
	if hashCodeVerifier != codeChallenge {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameCodeVerifier)
	}

	return nil
}

func (v *Validator) validateAudience() *oauth.Error {
	if err := v.validateAudienceRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateAudienceRegisteredForClient(); err != nil {
		return oauth.UnauthorizedClientError(err)
	}

	return nil
}

func (v *Validator) validateAudienceRequired() error {
	if v.params.Audience() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameAudience)
	}

	return nil
}

func (v *Validator) validateAudienceRegisteredForClient() error {
	if v.audienceRegisteredForClient() == false {
		return oauth.ErrOAuthAudienceNotRegisteredForClient
	}

	return nil
}

func (v *Validator) audienceRegisteredForClient() bool {
	if v.client.IsNone() {
		return false
	}

	client := v.client.Unwrap()
	return slices.Contains(client.Audiences(), v.params.Audience())
}
