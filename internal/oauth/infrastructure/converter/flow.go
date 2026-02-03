package converter

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/models"
)

func ToFlowRedis(flow dto.Flow) models.Flow {
	return models.Flow{
		ID:                  flow.ID().String(),
		ClientID:            flow.ClientID(),
		ResponseType:        flow.ResponseType(),
		RedirectURI:         flow.RedirectURI(),
		CodeChallenge:       flow.CodeChallenge(),
		CodeChallengeMethod: flow.CodeChallengeMethod(),
		State:               flow.State(),
		AuthorizationCode:   flow.AuthorizationCode(),
	}
}

func ToFlowDTO(flow models.Flow) (dto.Flow, error) {
	const op = "infrastructure.converter.ToFlowDTO"

	id, err := uuid.Parse(flow.ID)
	if err != nil {
		return dto.Flow{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.NewFlow(
		id,
		flow.ClientID,
		flow.ResponseType,
		flow.RedirectURI,
		flow.CodeChallenge,
		flow.CodeChallengeMethod,
		flow.State,
		dto.WithAuthorizationCode(flow.AuthorizationCode),
	), nil
}
