package converter

import (
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
)

func ToClientDTO(client models.Client, audiences []models.Audience) dto.Client {
	audienceURLs := make([]string, len(audiences))
	for i, audience := range audiences {
		audienceURLs[i] = audience.URL
	}

	return dto.Client{
		ID:        client.ID,
		Code:      client.Code,
		SecretKey: client.SecretKey,
		Audiences: audienceURLs,
		// TODO: get this from storage
		RedirectURI: []string{"http://localhost:3000"},
	}
}
