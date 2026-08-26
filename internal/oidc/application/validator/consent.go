package validator

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/pkg/extslices"
)

type consentRequestValidator struct {
	consentRequest   dto.ConsentRequest
	authorizeRequest dto.ValidatedAuthorizeRequest
}

func NewConsentRequestValidator(
	consentRequest dto.ConsentRequest,
	authorizeRequest dto.ValidatedAuthorizeRequest,
) *consentRequestValidator {
	return &consentRequestValidator{
		consentRequest:   consentRequest,
		authorizeRequest: authorizeRequest,
	}
}

func (v *consentRequestValidator) Validate() dto.ValidatedAuthorizeRequest {
	scopes := v.validateScopes()

	validatedData := dto.NewValidatedAuthorizeRequest(
		v.authorizeRequest.ResponseType(),
		v.authorizeRequest.Prompt(),
		v.authorizeRequest.ClientID(),
		v.authorizeRequest.RedirectURI(),
		v.authorizeRequest.CodeChallenge(),
		v.authorizeRequest.CodeChallengeMethod(),
		v.authorizeRequest.State(),
		v.authorizeRequest.Audience(),
		scopes,
	)

	return validatedData
}

func (v *consentRequestValidator) validateScopes() dto.ValidatedScopesRequest {
	grantedScopes := v.validateGrantedScopes()
	pendingScopes := v.validatePendingScopes()

	return dto.NewValidatedScopesRequest(dto.WithGranted(grantedScopes), dto.WithPending(pendingScopes))
}

func (v *consentRequestValidator) validateGrantedScopes() []string {
	grantedScopes := extslices.Intersect(v.authorizeRequest.AllScopes(), v.consentRequest.Scopes())
	return grantedScopes
}

func (v *consentRequestValidator) validatePendingScopes() []string {
	pendingScopes := extslices.Except(v.authorizeRequest.AllScopes(), v.consentRequest.Scopes())
	return pendingScopes
}
