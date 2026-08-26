package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache/models"
)

func ToAuthorizeRequestDTO(authorizationRequest models.AuthorizationRequest) dto.ValidatedAuthorizeRequest {
	scopes := dto.NewValidatedScopesRequest(
		dto.WithPending(authorizationRequest.Scopes.Pending),
		dto.WithGranted(authorizationRequest.Scopes.Granted),
	)

	return dto.NewValidatedAuthorizeRequest(
		authorizationRequest.ResponseType,
		authorizationRequest.Prompt,
		authorizationRequest.ClientID,
		authorizationRequest.RedirectURI,
		authorizationRequest.CodeChallenge,
		authorizationRequest.CodeChallengeMethod,
		authorizationRequest.State,
		authorizationRequest.Audience,
		scopes,
	)
}

func ToAuthorizationRequestStorage(authorizeRequest dto.ValidatedAuthorizeRequest) models.AuthorizationRequest {
	return models.AuthorizationRequest{
		ResponseType:        authorizeRequest.ResponseType(),
		Prompt:              authorizeRequest.Prompt(),
		ClientID:            authorizeRequest.ClientID(),
		RedirectURI:         authorizeRequest.RedirectURI(),
		CodeChallenge:       authorizeRequest.CodeChallenge(),
		CodeChallengeMethod: authorizeRequest.CodeChallengeMethod(),
		State:               authorizeRequest.State(),
		Audience:            authorizeRequest.Audience(),
		Scopes: models.ScopesRequest{
			Pending: authorizeRequest.PendingScopes(),
			Granted: authorizeRequest.GrantedScopes(),
		},
	}
}
