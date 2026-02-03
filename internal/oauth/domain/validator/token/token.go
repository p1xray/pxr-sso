package token

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/hasher/sha256"
	"github.com/p1xray/pxr-sso/pkg/nullable"
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

	return nil
}

func (v *Validator) validateFlowID() error {
	if err := v.validateFlowIDRequired(); err != nil {
		return err
	}

	if err := v.validateFlowIDValue(); err != nil {
		return err
	}

	if err := v.validateFlowIDExistFlow(); err != nil {
		return err
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

func (v *Validator) validateFlowIDExistFlow() error {
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

func (v *Validator) validateGrantType() error {
	if err := v.validateGrantTypeRequired(); err != nil {
		return err
	}

	if err := v.validateGrantTypeValue(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateGrantTypeRequired() error {
	if v.params.GrantType() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameGrantType)
	}

	return nil
}

func (v *Validator) validateGrantTypeValue() error {
	if v.params.GrantType() != oauth.RequestParameterAllowValueGrantTypeAuthorizationCode {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameGrantType)
	}

	return nil
}

func (v *Validator) validateClientID() error {
	if err := v.validateClientIDRequired(); err != nil {
		return err
	}

	if err := v.validateClientIDExistClient(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateClientIDRequired() error {
	if v.params.ClientID() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDExistClient() error {
	if v.client.IsNone() {
		return domain.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()

	if v.params.ClientID() != client.Code {
		return domain.ErrOAuthClientNotRegistered
	}

	return nil
}

func (v *Validator) validateAuthorizationCode() error {
	if err := v.validateAuthorizationCodeRequired(); err != nil {
		return err
	}

	if err := v.validateAuthorizationCodeEqualsFlowAuthorizationCode(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateAuthorizationCodeRequired() error {
	if v.params.AuthorizationCode() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameAuthorizationCode)
	}

	return nil
}

func (v *Validator) validateAuthorizationCodeEqualsFlowAuthorizationCode() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if v.params.AuthorizationCode() != flow.AuthorizationCode() {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameAuthorizationCode)
	}

	return nil
}

func (v *Validator) validateRedirectURI() error {
	if err := v.validateRedirectURIRequired(); err != nil {
		return err
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		return err
	}

	if err := v.validateRedirectURIEqualsFlowRedirectURI(); err != nil {
		return err
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

func (v *Validator) validateCodeVerifier() error {
	if err := v.validateCodeVerifierRequired(); err != nil {
		return err
	}

	if err := v.validateCodeVerifierHashEqualsFlowCodeChallenge(); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateCodeVerifierRequired() error {
	if v.params.CodeVerifier() == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameCodeVerifier)
	}

	return nil
}

func (v *Validator) validateCodeVerifierHashEqualsFlowCodeChallenge() error {
	if v.flow.IsNone() {
		return domain.ErrOAuthFlowNotExists
	}

	flow := v.flow.Unwrap()
	if sha256.Base64URLHash(v.params.CodeVerifier()) != flow.CodeChallenge() {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameCodeVerifier)
	}

	return nil
}
