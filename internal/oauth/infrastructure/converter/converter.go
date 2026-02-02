package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	redisModels "github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/models"
	storageModels "github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

func ToClientDTO(client storageModels.Client, audiences []storageModels.Audience) dto.Client {
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

func ToFlowRedis(flow dto.Flow) redisModels.Flow {
	return redisModels.Flow{
		ID:                  flow.ID().String(),
		ClientID:            flow.ClientID(),
		RedirectURI:         flow.RedirectURI(),
		CodeChallenge:       flow.CodeChallenge(),
		CodeChallengeMethod: flow.CodeChallengeMethod(),
		State:               flow.State(),
		AuthorizationCode:   flow.AuthorizationCode(),
	}
}
