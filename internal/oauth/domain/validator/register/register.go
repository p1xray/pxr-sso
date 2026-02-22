package register

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
	params dto.Register
	client nullable.Nullable[dto.Client]
	flow   nullable.Nullable[dto.Flow]
	user   nullable.Nullable[dto.User]
}

func NewValidator(
	params dto.Register,
	client nullable.Nullable[dto.Client],
	flow nullable.Nullable[dto.Flow],
	user nullable.Nullable[dto.User],
) *Validator {
	return &Validator{
		params: params,
		client: client,
		flow:   flow,
		user:   user,
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

	if err := v.validateUsername(); err != nil {
		return err
	}

	if err := v.validatePassword(); err != nil {
		return err
	}

	if err := v.validateFullName(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateFlowID() *oauth.Error {
	if err := v.validateFlowIDRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateFlowIDValue(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateFlowIDExistsFlow(); err != nil {
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

func (v *Validator) validateFlowIDExistsFlow() error {
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

func (v *Validator) validateResponseType() *oauth.Error {
	if err := v.validateResponseTypeRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateResponseTypeValue(); err != nil {
		return oauth.UnsupportedResponseTypeError(err)
	}

	if err := v.validateResponseTypeEqualsFlowResponseType(); err != nil {
		return oauth.UnsupportedResponseTypeError(err)
	}

	return nil
}

func (v *Validator) validateResponseTypeRequired() error {
	if v.params.ResponseType() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeValue() error {
	if v.params.ResponseType() != validator.RequestParameterAllowValueResponseTypeCode {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeEqualsFlowResponseType() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.ResponseType() != flow.ResponseType() {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateClientID() *oauth.Error {
	if err := v.validateClientIDRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateClientIDExistsClient(); err != nil {
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

func (v *Validator) validateClientIDExistsClient() error {
	if v.client.IsNone() {
		return oauth.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()

	if v.params.ClientID() != client.Code() {
		return oauth.ErrOAuthClientNotRegistered
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

func (v *Validator) validateState() *oauth.Error {
	if err := v.validateStateRequired(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	if err := v.validateStateEqualsFlowState(); err != nil {
		return oauth.InvalidRequestError(err)
	}

	return nil
}

func (v *Validator) validateStateRequired() error {
	if v.params.State() == "" {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateStateEqualsFlowState() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.State() != flow.State() {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateScope() *oauth.Error {
	if err := v.validateScopeEqualsFlowScope(); err != nil {
		return oauth.InvalidScopeError(err)
	}

	return nil
}

func (v *Validator) validateScopeEqualsFlowScope() error {
	if v.flow.IsNone() {
		return oauth.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	equals := extslices.Any(flow.Scope(), v.params.Scope())
	if !equals {
		return fmt.Errorf("%w: %s", oauth.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameScope)
	}

	return nil
}

func (v *Validator) validateUsername() *oauth.Error {
	if err := v.validateUsernameRequired(); err != nil {
		return oauth.InvalidUserCredentialsError(err)
	}

	if err := v.validateUsernameUserNotExists(); err != nil {
		return oauth.InvalidUserCredentialsError(err)
	}

	return nil
}

func (v *Validator) validateUsernameRequired() error {
	if v.params.Username() == "" {
		return oauth.ErrOAuthUsernameRequired
	}

	return nil
}

func (v *Validator) validateUsernameUserNotExists() error {
	if v.user.IsNone() {
		return nil
	}

	user := v.user.Unwrap()
	if v.params.Username() == user.Username() {
		return oauth.ErrOAuthUserAlreadyExists
	}

	return nil
}

func (v *Validator) validatePassword() *oauth.Error {
	if err := v.validatePasswordRequired(); err != nil {
		return oauth.InvalidUserCredentialsError(err)
	}

	return nil
}

func (v *Validator) validatePasswordRequired() error {
	if v.params.Password() == "" {
		return oauth.ErrOAuthPasswordRequired
	}

	return nil
}

func (v *Validator) validateFullName() *oauth.Error {
	if err := v.validateFullNameRequired(); err != nil {
		return oauth.InvalidUserCredentialsError(err)
	}

	return nil
}

func (v *Validator) validateFullNameRequired() error {
	if v.params.FullName() == "" {
		return oauth.ErrOAuthFullNameRequired
	}

	return nil
}
