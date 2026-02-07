package register

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/nullable"
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

func (v *Validator) Validate() *domain.DisplayableError {
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

func (v *Validator) validateFlowID() *domain.DisplayableError {
	if err := v.validateFlowIDRequired(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateFlowIDValue(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateFlowIDExistsFlow(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (v *Validator) validateFlowIDRequired() error {
	if v.params.FlowID() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDValue() error {
	if err := uuid.Validate(v.params.FlowID()); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	return nil
}

func (v *Validator) validateFlowIDExistsFlow() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()

	paramFlowID, err := uuid.Parse(v.params.FlowID())
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameFlowID)
	}

	if paramFlowID != flow.ID() {
		return domain.ErrOAuthFlowNotExists
	}

	return nil
}

func (v *Validator) validateResponseType() *domain.DisplayableError {
	if err := v.validateResponseTypeRequired(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateResponseTypeValue(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateResponseTypeEqualsFlowResponseType(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (v *Validator) validateResponseTypeRequired() error {
	if v.params.ResponseType() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeValue() error {
	if v.params.ResponseType() != oauth.RequestParameterAllowValueResponseTypeCode {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeEqualsFlowResponseType() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.ResponseType() != flow.ResponseType() {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateClientID() *domain.DisplayableError {
	if err := v.validateClientIDRequired(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateClientIDExistsClient(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (v *Validator) validateClientIDRequired() error {
	if v.params.ClientID() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDExistsClient() error {
	if v.client.IsNone() {
		return domain.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()

	if v.params.ClientID() != client.Code {
		return domain.ErrOAuthClientNotRegistered
	}

	return nil
}

func (v *Validator) validateRedirectURI() *domain.DisplayableError {
	if err := v.validateRedirectURIRequired(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateRedirectURIEqualsFlowRedirectURI(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (v *Validator) validateRedirectURIRequired() error {
	if v.params.RedirectURI() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIEqualsFlowRedirectURI() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.RedirectURI() != flow.RedirectURI() {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIRegisteredForClient() error {
	if v.redirectURIRegisteredForClient() == false {
		return domain.ErrOAuthRedirectURINotRegisteredForClient
	}

	return nil
}

func (v *Validator) redirectURIRegisteredForClient() bool {
	if v.client.IsNone() {
		return false
	}

	client := v.client.Unwrap()
	for _, clientRedirectURI := range client.RedirectURI {
		if v.params.RedirectURI() == clientRedirectURI {
			return true
		}
	}

	return false
}

func (v *Validator) validateState() *domain.DisplayableError {
	if err := v.validateStateRequired(); err != nil {
		return domain.InternalError(err)
	}

	if err := v.validateStateEqualsFlowState(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (v *Validator) validateStateRequired() error {
	if v.params.State() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateStateEqualsFlowState() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.State() != flow.State() {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateUsername() *domain.DisplayableError {
	if err := v.validateUsernameRequired(); err != nil {
		return domain.DisplayError(domain.ErrorDescriptionUsernameRequired, err)
	}

	if err := v.validateUsernameUserNotExists(); err != nil {
		return domain.DisplayError(domain.ErrorDescriptionInvalidUserCredentials, err)
	}

	return nil
}

func (v *Validator) validateUsernameRequired() error {
	if v.params.Username() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameUsername)
	}

	return nil
}

func (v *Validator) validateUsernameUserNotExists() error {
	if v.user.IsNone() {
		return nil
	}

	user := v.user.Unwrap()
	if v.params.Username() == user.Username() {
		return domain.ErrOAuthUserAlreadyExists
	}

	return nil
}

func (v *Validator) validatePassword() *domain.DisplayableError {
	if err := v.validatePasswordRequired(); err != nil {
		return domain.DisplayError(domain.ErrorDescriptionPasswordRequired, err)
	}

	return nil
}

func (v *Validator) validatePasswordRequired() error {
	if v.params.Password() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNamePassword)
	}

	return nil
}

func (v *Validator) validateFullName() *domain.DisplayableError {
	if err := v.validateFullNameRequired(); err != nil {
		return domain.DisplayError(domain.ErrorDescriptionUsernameRequired, err)
	}

	return nil
}

func (v *Validator) validateFullNameRequired() error {
	if v.params.FullName() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameFullName)
	}

	return nil
}
