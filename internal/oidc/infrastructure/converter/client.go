package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToClientDTO(client models.Client) dto.Client {
	audienceURIs := make([]string, len(client.Audiences))
	for i, audience := range client.Audiences {
		audienceURIs[i] = audience.URI
	}

	redirectURIs := make([]string, len(client.RedirectURIs))
	for i, redirectURI := range client.RedirectURIs {
		redirectURIs[i] = redirectURI.URI
	}

	scopes := make([]string, len(client.ScopeLinks))
	for i, link := range client.ScopeLinks {
		scopes[i] = link.Scope.Code
	}

	defaultRoles := make([]dto.Role, len(client.DefaultRoleLinks))
	for i, link := range client.DefaultRoleLinks {
		roleDTO := ToRoleDTO(link.Role)
		defaultRoles[i] = roleDTO
	}

	clientDTO := dto.NewClient(
		client.ID,
		client.Code,
		client.SecretKey,
		audienceURIs,
		redirectURIs,
		scopes,
		defaultRoles,
	)

	return clientDTO
}
