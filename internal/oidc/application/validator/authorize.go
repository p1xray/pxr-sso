package validator

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"slices"
)

type authorizationRequestValidator struct {
	request  dto.AuthorizeRequest
	client   nullable.Nullable[dto.Client]
	sessions []dto.Session

	err *oidc.Error

	isResponseTypeValid        bool
	isPromptValid              bool
	isClientIDValid            bool
	isRedirectURIValid         bool
	isCodeChallengeValid       bool
	isCodeChallengeMethodValid bool
	isStateValid               bool
	isAudienceValid            bool
}

func NewAuthorizationRequestValidator(
	request dto.AuthorizeRequest,
	client nullable.Nullable[dto.Client],
	sessions []dto.Session,
) *authorizationRequestValidator {
	return &authorizationRequestValidator{
		request:  request,
		client:   client,
		sessions: sessions,
	}
}

func (v *authorizationRequestValidator) Validate() (dto.ValidatedAuthorizeRequest, error) {
	v.validateResponseType()
	v.validatePrompt()
	v.validateClientID()
	v.validateRedirectURI()
	v.validateCodeChallenge()
	v.validateCodeChallengeMethod()
	v.validateState()
	v.validateAudience()

	validatedRequest := v.validatedRequest()

	return validatedRequest, v.err
}

func (v *authorizationRequestValidator) validatedRequest() dto.ValidatedAuthorizeRequest {
	responseType := v.validatedResponseType()
	prompt := v.validatedPrompt()
	clientID := v.validatedClientID()
	redirectURI := v.validatedRedirectURI()
	codeChallenge := v.validatedCodeChallenge()
	codeChallengeMethod := v.validatedCodeChallengeMethod()
	state := v.validatedState()
	audience := v.validatedAudience()
	scopes := v.validatedScope()

	validatedData := dto.NewValidatedAuthorizeRequest(
		responseType,
		prompt,
		clientID,
		redirectURI,
		codeChallenge,
		codeChallengeMethod,
		state,
		audience,
		scopes,
	)

	return validatedData
}

func (v *authorizationRequestValidator) validateResponseType() {
	if err := v.validateResponseTypeRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateResponseTypeOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateResponseTypeValue(); err != nil {
		v.setErrorIfEmpty(oidc.UnsupportedResponseTypeError(err))
	}

	v.isResponseTypeValid = true
}

func (v *authorizationRequestValidator) validateResponseTypeRequired() error {
	if len(v.request.ResponseType()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameResponseType)
	}

	return nil
}

func (v *authorizationRequestValidator) validateResponseTypeOnlyOneValue() error {
	if len(v.request.ResponseType()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameResponseType)
	}

	return nil
}

func (v *authorizationRequestValidator) validateResponseTypeValue() error {
	paramResponseType := v.request.ResponseType()
	if len(paramResponseType) == 1 && paramResponseType[0] != RequestParameterAllowValueResponseTypeCode {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameResponseType)
	}

	return nil
}

func (v *authorizationRequestValidator) validatedResponseType() string {
	if !v.isResponseTypeValid {
		return ""
	}

	paramResponseType := v.request.ResponseType()
	responseType := paramResponseType[0]

	return responseType
}

func (v *authorizationRequestValidator) validatePrompt() {
	if err := v.validatePromptRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validatePromptOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validatePromptValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	v.isPromptValid = true
}

func (v *authorizationRequestValidator) validatePromptRequired() error {
	if len(v.request.Prompt()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNamePrompt)
	}

	return nil
}

func (v *authorizationRequestValidator) validatePromptOnlyOneValue() error {
	if len(v.request.Prompt()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNamePrompt)
	}

	return nil
}

func (v *authorizationRequestValidator) validatePromptValue() error {
	paramPrompt := v.request.Prompt()
	if len(paramPrompt) == 1 && slices.Contains(RequestParameterAllowValuePrompt, paramPrompt[0]) == false {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNamePrompt)
	}

	return nil
}

func (v *authorizationRequestValidator) validatedPrompt() string {
	if !v.isPromptValid {
		return ""
	}

	paramPrompt := v.request.Prompt()
	prompt := paramPrompt[0]

	return prompt
}

func (v *authorizationRequestValidator) validateClientID() {
	if err := v.validateClientIDRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateClientIDOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateClientIDExistClient(); err != nil {
		v.setErrorIfEmpty(oidc.UnauthorizedClientError(err))
	}

	v.isClientIDValid = true
}

func (v *authorizationRequestValidator) validateClientIDRequired() error {
	if len(v.request.ClientID()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameClientID)
	}

	return nil
}

func (v *authorizationRequestValidator) validateClientIDOnlyOneValue() error {
	if len(v.request.ClientID()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameClientID)
	}

	return nil
}

func (v *authorizationRequestValidator) validateClientIDExistClient() error {
	if v.client.IsNone() {
		return oidc.ErrOAuthClientNotRegistered
	}

	client := v.client.Unwrap()
	paramClientID := v.request.ClientID()
	if len(paramClientID) == 1 && paramClientID[0] == client.Code() {
		return nil
	}

	return oidc.ErrOAuthClientNotRegistered
}

func (v *authorizationRequestValidator) validatedClientID() string {
	if !v.isClientIDValid {
		return ""
	}

	paramClientID := v.request.ClientID()
	clientID := paramClientID[0]

	return clientID
}

func (v *authorizationRequestValidator) validateRedirectURI() {
	if err := v.validateRedirectURIRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateRedirectURIOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateRedirectURIRegisteredForClient(); err != nil {
		v.setErrorIfEmpty(oidc.UnauthorizedClientError(err))
	}

	v.isRedirectURIValid = true
}

func (v *authorizationRequestValidator) validateRedirectURIRequired() error {
	if len(v.request.RedirectURI()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *authorizationRequestValidator) validateRedirectURIOnlyOneValue() error {
	if len(v.request.RedirectURI()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameRedirectURI)
	}

	return nil
}

func (v *authorizationRequestValidator) validateRedirectURIRegisteredForClient() error {
	if v.redirectURIRegisteredForClient() == false {
		return oidc.ErrOAuthRedirectURINotRegisteredForClient
	}

	return nil
}

func (v *authorizationRequestValidator) redirectURIRegisteredForClient() bool {
	if v.client.IsNone() {
		return false
	}

	client := v.client.Unwrap()
	return extslices.Any(client.RedirectURI(), v.request.RedirectURI())
}

func (v *authorizationRequestValidator) validatedRedirectURI() string {
	if !v.isRedirectURIValid {
		return ""
	}

	paramRedirectURI := v.request.RedirectURI()
	redirectURI := paramRedirectURI[0]

	return redirectURI
}

func (v *authorizationRequestValidator) validateCodeChallenge() {
	if err := v.validateCodeChallengeRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateCodeChallengeOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateCodeChallengeValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	v.isCodeChallengeValid = true
}

func (v *authorizationRequestValidator) validateCodeChallengeRequired() error {
	if len(v.request.CodeChallenge()) == 0 {
		return fmt.Errorf(
			"%w: %s",
			oidc.ErrOAuthMissingRequiredParameter,
			oidc.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *authorizationRequestValidator) validateCodeChallengeOnlyOneValue() error {
	if len(v.request.CodeChallenge()) > 1 {
		return fmt.Errorf(
			"%w: %s",
			oidc.ErrOAuthParameterIncludedMoreThanOnce,
			oidc.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *authorizationRequestValidator) validateCodeChallengeValue() error {
	paramCodeChallengeMethod := v.request.CodeChallenge()
	if len(paramCodeChallengeMethod) == 1 && paramCodeChallengeMethod[0] == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameCodeChallenge)
	}

	return nil
}

func (v *authorizationRequestValidator) validatedCodeChallenge() string {
	if !v.isCodeChallengeValid {
		return ""
	}

	paramCodeChallenge := v.request.CodeChallenge()
	codeChallenge := paramCodeChallenge[0]

	return codeChallenge
}

func (v *authorizationRequestValidator) validateCodeChallengeMethod() {
	if err := v.validateCodeChallengeMethodRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateCodeChallengeMethodOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateCodeChallengeMethodValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	v.isCodeChallengeMethodValid = true
}

func (v *authorizationRequestValidator) validateCodeChallengeMethodRequired() error {
	if len(v.request.CodeChallengeMethod()) == 0 {
		return fmt.Errorf(
			"%w: %s",
			oidc.ErrOAuthMissingRequiredParameter,
			oidc.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *authorizationRequestValidator) validateCodeChallengeMethodOnlyOneValue() error {
	if len(v.request.CodeChallengeMethod()) > 1 {
		return fmt.Errorf(
			"%w: %s",
			oidc.ErrOAuthParameterIncludedMoreThanOnce,
			oidc.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *authorizationRequestValidator) validateCodeChallengeMethodValue() error {
	paramCodeChallengeMethod := v.request.CodeChallengeMethod()
	if len(paramCodeChallengeMethod) == 1 &&
		paramCodeChallengeMethod[0] != RequestParameterAllowValueCodeChallengeMethod {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameCodeChallengeMethod)
	}

	return nil
}

func (v *authorizationRequestValidator) validatedCodeChallengeMethod() string {
	if !v.isCodeChallengeMethodValid {
		return ""
	}

	paramCodeChallengeMethod := v.request.CodeChallengeMethod()
	codeChallengeMethod := paramCodeChallengeMethod[0]

	return codeChallengeMethod
}

func (v *authorizationRequestValidator) validateState() {
	if err := v.validateStateRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateStateOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateStateValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	v.isStateValid = true
}

func (v *authorizationRequestValidator) validateStateRequired() error {
	if len(v.request.State()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameState)
	}

	return nil
}

func (v *authorizationRequestValidator) validateStateOnlyOneValue() error {
	if len(v.request.State()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameState)
	}

	return nil
}

func (v *authorizationRequestValidator) validateStateValue() error {
	paramState := v.request.State()
	if len(paramState) == 1 && paramState[0] == "" {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterInvalidValue, oidc.RequestParameterNameState)
	}

	return nil
}

func (v *authorizationRequestValidator) validatedState() string {
	if !v.isStateValid {
		return ""
	}

	paramState := v.request.State()
	state := paramState[0]

	return state
}

func (v *authorizationRequestValidator) validateAudience() {
	if err := v.validateAudienceRequired(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateAudienceOnlyOneValue(); err != nil {
		v.setErrorIfEmpty(oidc.InvalidRequestError(err))
	}

	if err := v.validateAudienceRegisteredForClient(); err != nil {
		v.setErrorIfEmpty(oidc.UnauthorizedClientError(err))
	}

	v.isAudienceValid = true
}

func (v *authorizationRequestValidator) validateAudienceRequired() error {
	if len(v.request.Audience()) == 0 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthMissingRequiredParameter, oidc.RequestParameterNameAudience)
	}

	return nil
}

func (v *authorizationRequestValidator) validateAudienceOnlyOneValue() error {
	if len(v.request.Audience()) > 1 {
		return fmt.Errorf("%w: %s", oidc.ErrOAuthParameterIncludedMoreThanOnce, oidc.RequestParameterNameAudience)
	}

	return nil
}

func (v *authorizationRequestValidator) validateAudienceRegisteredForClient() error {
	if v.audienceRegisteredForClient() == false {
		return oidc.ErrOAuthAudienceNotRegisteredForClient
	}

	return nil
}

func (v *authorizationRequestValidator) audienceRegisteredForClient() bool {
	if v.client.IsNone() {
		return false
	}

	client := v.client.Unwrap()

	paramAudience := v.request.Audience()
	if len(paramAudience) == 0 {
		return false
	}

	return slices.Contains(client.Audiences(), paramAudience[0])
}

func (v *authorizationRequestValidator) validatedAudience() string {
	if !v.isAudienceValid {
		return ""
	}

	paramAudience := v.request.Audience()
	audience := paramAudience[0]

	return audience
}

func (v *authorizationRequestValidator) validatedScope() dto.ValidatedScopesRequest {
	if len(v.sessions) > 1 {
		return dto.NewValidatedScopesRequest()
	}

	if v.client.IsNone() {
		return dto.NewValidatedScopesRequest()
	}

	client := v.client.Unwrap()
	scopesValidatedByClient := extslices.Intersect(client.AvailableScopes(), v.request.Scope())

	if len(v.sessions) == 0 {
		validateScopes := dto.NewValidatedScopesRequest(dto.WithPending(scopesValidatedByClient))
		return validateScopes
	}

	session := v.sessions[0]

	pendingScopes := extslices.Except(scopesValidatedByClient, session.ScopeCodes())
	grantedScopes := extslices.Intersect(scopesValidatedByClient, session.ScopeCodes())
	validatedScopes := dto.NewValidatedScopesRequest(dto.WithPending(pendingScopes), dto.WithGranted(grantedScopes))

	return validatedScopes
}

func (v *authorizationRequestValidator) setErrorIfEmpty(err *oidc.Error) {
	if v.err == nil {
		v.err = err
	}
}
