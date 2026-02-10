package authorize

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/nullable"
)

type Validator struct {
	params dto.Authorize
	client nullable.Nullable[dto.Client]

	err *domain.OAuthError

	isValidatingExecuted     bool
	responseTypeValid        bool
	clientIDValid            bool
	redirectURIValid         bool
	codeChallengeValid       bool
	codeChallengeMethodValid bool
	stateValid               bool
}

func NewValidator(
	params dto.Authorize,
	client nullable.Nullable[dto.Client],
) *Validator {
	return &Validator{params: params, client: client}
}

func (v *Validator) Validate() *domain.OAuthError {
	v.isValidatingExecuted = true

	v.validateResponseType()
	v.validateClientID()
	v.validateRedirectURI()
	v.validateCodeChallenge()
	v.validateCodeChallengeMethod()
	v.validateState()

	return v.err
}

func (v *Validator) ValidatedData() dto.ValidatedAuthorize {
	if !v.isValidatingExecuted {
		_ = v.Validate()
	}

	responseType := v.validatedResponseType()
	clientID := v.validatedClientID()
	redirectURI := v.validatedRedirectURI()
	codeChallenge := v.validatedCodeChallenge()
	codeChallengeMethod := v.validatedCodeChallengeMethod()
	state := v.validatedState()
	scope := v.validatedScope()

	validatedData := dto.NewValidatedAuthorize(
		responseType,
		clientID,
		redirectURI,
		codeChallenge,
		codeChallengeMethod,
		state,
		scope,
	)

	return validatedData
}

func (v *Validator) validateResponseType() {
	if err := v.validateResponseTypeRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateResponseTypeMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateResponseTypeValue(); err != nil {
		v.setErrorIfEmpty(domain.UnsupportedResponseTypeOAuthError(err))
	}

	v.responseTypeValid = true
}

func (v *Validator) validateResponseTypeRequired() error {
	if len(v.params.ResponseType()) == 0 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeMoreThenOnce() error {
	if len(v.params.ResponseType()) > 1 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterIncludedMoreThanOnce, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) validateResponseTypeValue() error {
	paramResponseType := v.params.ResponseType()
	if len(paramResponseType) == 1 && paramResponseType[0] != oauth.RequestParameterAllowValueResponseTypeCode {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameResponseType)
	}

	return nil
}

func (v *Validator) isResponseTypeValid() bool {
	if !v.isValidatingExecuted {
		v.validatedResponseType()
	}

	return v.responseTypeValid
}

func (v *Validator) validatedResponseType() string {
	if !v.isResponseTypeValid() {
		return ""
	}

	paramResponseType := v.params.ResponseType()
	responseType := paramResponseType[0]

	return responseType
}

func (v *Validator) validateClientID() {
	if err := v.validateClientIDRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateClientIDMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateClientIDExistClient(); err != nil {
		v.setErrorIfEmpty(domain.UnauthorizedClientOAuthError(err))
	}

	v.clientIDValid = true
}

func (v *Validator) validateClientIDRequired() error {
	if len(v.params.ClientID()) == 0 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDMoreThenOnce() error {
	if len(v.params.ClientID()) > 1 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterIncludedMoreThanOnce, oauth.RequestParameterNameClientID)
	}

	return nil
}

func (v *Validator) validateClientIDExistClient() error {
	if v.client.IsNone() {
		return domain.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()
	paramClientID := v.params.ClientID()
	if len(paramClientID) == 1 && paramClientID[0] == client.Code() {
		return nil
	}

	return domain.ErrOAuthClientNotRegistered
}

func (v *Validator) isClientIDValid() bool {
	if !v.isValidatingExecuted {
		v.validateClientID()
	}

	return v.clientIDValid
}

func (v *Validator) validatedClientID() string {
	if !v.isClientIDValid() {
		return ""
	}

	paramClientID := v.params.ClientID()
	clientID := paramClientID[0]

	return clientID
}

func (v *Validator) validateRedirectURI() {
	if err := v.validateRedirectURIRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateRedirectURIMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		v.setErrorIfEmpty(domain.UnauthorizedClientOAuthError(err))
	}

	v.redirectURIValid = true
}

func (v *Validator) validateRedirectURIRequired() error {
	if len(v.params.RedirectURI()) == 0 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *Validator) validateRedirectURIMoreThenOnce() error {
	if len(v.params.RedirectURI()) > 1 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterIncludedMoreThanOnce, oauth.RequestParameterNameRedirectURI)
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
	return extslices.Any(client.RedirectURI(), v.params.RedirectURI())
}

func (v *Validator) isRedirectURIValid() bool {
	if !v.isValidatingExecuted {
		v.validateRedirectURI()
	}

	return v.redirectURIValid
}

func (v *Validator) validatedRedirectURI() string {
	if !v.isRedirectURIValid() {
		return ""
	}

	paramRedirectURI := v.params.RedirectURI()
	redirectURI := paramRedirectURI[0]

	return redirectURI
}

func (v *Validator) validateCodeChallenge() {
	if err := v.validateCodeChallengeRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateCodeChallengeMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateCodeChallengeValue(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	v.codeChallengeValid = true
}

func (v *Validator) validateCodeChallengeRequired() error {
	if len(v.params.CodeChallenge()) == 0 {
		return fmt.Errorf(
			"%w: %s",
			domain.ErrOAuthMissingRequiredParameter,
			oauth.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *Validator) validateCodeChallengeMoreThenOnce() error {
	if len(v.params.CodeChallenge()) > 1 {
		return fmt.Errorf(
			"%w: %s",
			domain.ErrOAuthParameterIncludedMoreThanOnce,
			oauth.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *Validator) validateCodeChallengeValue() error {
	paramCodeChallengeMethod := v.params.CodeChallenge()
	if len(paramCodeChallengeMethod) == 1 && paramCodeChallengeMethod[0] == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *Validator) isCodeChallengeValid() bool {
	if !v.isValidatingExecuted {
		v.validateCodeChallenge()
	}

	return v.codeChallengeValid
}

func (v *Validator) validatedCodeChallenge() string {
	if !v.isCodeChallengeValid() {
		return ""
	}

	paramCodeChallenge := v.params.CodeChallenge()
	codeChallenge := paramCodeChallenge[0]

	return codeChallenge
}

func (v *Validator) validateCodeChallengeMethod() {
	if err := v.validateCodeChallengeMethodRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateCodeChallengeMethodMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateCodeChallengeMethodValue(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	v.codeChallengeMethodValid = true
}

func (v *Validator) validateCodeChallengeMethodRequired() error {
	if len(v.params.CodeChallengeMethod()) == 0 {
		return fmt.Errorf(
			"%w: %s",
			domain.ErrOAuthMissingRequiredParameter,
			oauth.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *Validator) validateCodeChallengeMethodMoreThenOnce() error {
	if len(v.params.CodeChallengeMethod()) > 1 {
		return fmt.Errorf(
			"%w: %s",
			domain.ErrOAuthParameterIncludedMoreThanOnce,
			oauth.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *Validator) validateCodeChallengeMethodValue() error {
	paramCodeChallengeMethod := v.params.CodeChallengeMethod()
	if len(paramCodeChallengeMethod) == 1 &&
		paramCodeChallengeMethod[0] != oauth.RequestParameterAllowValueCodeChallengeMethod {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *Validator) isCodeChallengeMethodValid() bool {
	if !v.isValidatingExecuted {
		v.validateCodeChallengeMethod()
	}

	return v.codeChallengeMethodValid
}

func (v *Validator) validatedCodeChallengeMethod() string {
	if !v.isCodeChallengeMethodValid() {
		return ""
	}

	paramCodeChallengeMethod := v.params.CodeChallengeMethod()
	codeChallengeMethod := paramCodeChallengeMethod[0]

	return codeChallengeMethod
}

func (v *Validator) validateState() {
	if err := v.validateStateRequired(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateStateMoreThenOnce(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	if err := v.validateStateValue(); err != nil {
		v.setErrorIfEmpty(domain.InvalidRequestOAuthError(err))
	}

	v.stateValid = true
}

func (v *Validator) validateStateRequired() error {
	if len(v.params.State()) == 0 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthMissingRequiredParameter, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateStateMoreThenOnce() error {
	if len(v.params.State()) > 1 {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterIncludedMoreThanOnce, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) validateStateValue() error {
	paramState := v.params.State()
	if len(paramState) == 1 && paramState[0] == "" {
		return fmt.Errorf("%w: %s", domain.ErrOAuthParameterInvalidValue, oauth.RequestParameterNameState)
	}

	return nil
}

func (v *Validator) isStateValid() bool {
	if !v.isValidatingExecuted {
		v.validateState()
	}

	return v.stateValid
}

func (v *Validator) validatedState() string {
	if !v.isStateValid() {
		return ""
	}

	paramState := v.params.State()
	state := paramState[0]

	return state
}

func (v *Validator) validatedScope() []string {
	if v.client.IsNone() {
		return []string{}
	}

	client := v.client.Unwrap()
	scope := extslices.Intersect(client.Scope(), v.params.Scope())

	return scope
}

func (v *Validator) setErrorIfEmpty(err *domain.OAuthError) {
	if v.err == nil {
		v.err = err
	}
}
