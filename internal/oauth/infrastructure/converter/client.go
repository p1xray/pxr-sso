package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

func ToClientDTO(
	client models.Client,
	clientAudiences []models.Audience,
	clientRedirectURIs []models.RedirectURI,
	clientScopes []models.Scope,
) dto.Client {
	audienceURIs := make([]string, len(clientAudiences))
	for i, audience := range clientAudiences {
		audienceURIs[i] = audience.URI
	}

	redirectURIs := make([]string, len(clientRedirectURIs))
	for i, redirectURI := range clientRedirectURIs {
		redirectURIs[i] = redirectURI.URI
	}

	scopes := make([]string, len(clientScopes))
	for i, scope := range clientScopes {
		scopes[i] = scope.Code
	}

	return dto.Client{
		ID:          client.ID,
		Code:        client.Code,
		SecretKey:   client.SecretKey,
		Audiences:   audienceURIs,
		RedirectURI: redirectURIs,
		Scope:       scopes,
	}
}
