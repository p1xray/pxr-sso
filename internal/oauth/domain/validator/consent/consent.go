package consent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"slices"
)

type Validator struct {
	params dto.Consent
	client nullable.Nullable[dto.Client]
	flow   nullable.Nullable[dto.Flow]
}

func NewValidator(
	params dto.Consent,
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

	if err := v.validateResponseType(); err != nil {
		return err
	}

	if err := v.validateClientID(); err != nil {
		return err
	}

	if err := v.validateRedirectURI(); err != nil {
		return err
	}

	if err := v.validateState(); err != nil {
		return err
	}

	if err := v.validateScope(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateFlowID() *validator.Error {
	if err := v.validateFlowIDRequired(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateFlowIDValue(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateFlowIDExistsFlow(); err != nil {
		return validator.InvalidRequestError(err)
	}

	return nil
}

func (v *Validator) validateFlowIDRequired() error {
	if v.params.FlowID() == "" {
		return fmt.Errorf("%w: %s", validator.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDValue() error {
	if err := uuid.Validate(v.params.FlowID()); err != nil {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDExistsFlow() error {
	if v.flow.IsNone() {
		return validator.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()

	paramFlowID, err := uuid.Parse(v.params.FlowID())
	if err != nil {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	if paramFlowID != flow.ID() {
		return validator.ErrOAuthFlowNotExists
	}

	return nil
}

func (v *Validator) validateResponseType() *validator.Error {
	if err := v.validateResponseTypeRequired(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateResponseTypeValue(); err != nil {
		return validator.UnsupportedResponseTypeError(err)
	}

	if err := v.validateResponseTypeEqualsFlowResponseType(); err != nil {
		return validator.UnsupportedResponseTypeError(err)
	}

	return nil
}

func (v *Validator) validateResponseTypeRequired() error {
	if v.params.ResponseType() == "" {
		return fmt.Errorf("%w: %s", validator.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeValue() error {
	if v.params.ResponseType() != validator.RequestParameterAllowValueResponseTypeCode {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeEqualsFlowResponseType() error {
	if v.flow.IsNone() {
		return validator.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.ResponseType() != flow.ResponseType() {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateClientID() *validator.Error {
	if err := v.validateClientIDRequired(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateClientIDExistsClient(); err != nil {
		return validator.UnauthorizedClientError(err)
	}

	return nil
}

func (v *Validator) validateClientIDRequired() error {
	if v.params.ClientID() == "" {
		return fmt.Errorf("%w: %s", validator.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDExistsClient() error {
	if v.client.IsNone() {
		return validator.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()

	if v.params.ClientID() != client.Code() {
		return validator.ErrOAuthClientNotRegistered
	}

	return nil
}

func (v *Validator) validateRedirectURI() *validator.Error {
	if err := v.validateRedirectURIRequired(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		return validator.UnauthorizedClientError(err)
	}

	if err := v.validateRedirectURIEqualsFlowRedirectURI(); err != nil {
		return validator.UnauthorizedClientError(err)
	}

	return nil
}

func (v *Validator) validateRedirectURIRequired() error {
	if v.params.RedirectURI() == "" {
		return fmt.Errorf("%w: %s", validator.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIEqualsFlowRedirectURI() error {
	if v.flow.IsNone() {
		return validator.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.RedirectURI() != flow.RedirectURI() {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIRegisteredForClient() error {
	if v.redirectURIRegisteredForClient() == false {
		return validator.ErrOAuthRedirectURINotRegisteredForClient
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

func (v *Validator) validateState() *validator.Error {
	if err := v.validateStateRequired(); err != nil {
		return validator.InvalidRequestError(err)
	}

	if err := v.validateStateEqualsFlowState(); err != nil {
		return validator.InvalidRequestError(err)
	}

	return nil
}

func (v *Validator) validateStateRequired() error {
	if v.params.State() == "" {
		return fmt.Errorf("%w: %s", validator.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateStateEqualsFlowState() error {
	if v.flow.IsNone() {
		return validator.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.State() != flow.State() {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateScope() *validator.Error {
	if err := v.validateScopeEqualsFlowScope(); err != nil {
		return validator.InvalidScopeError(err)
	}

	return nil
}

func (v *Validator) validateScopeEqualsFlowScope() error {
	if v.flow.IsNone() {
		return validator.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	equals := extslices.Any(flow.Scope(), v.params.Scope())
	if !equals {
		return fmt.Errorf("%w: %s", validator.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameScope)
	}

	return nil
}
