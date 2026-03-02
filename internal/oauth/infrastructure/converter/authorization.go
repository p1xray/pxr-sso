package converter

import (
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/models"
	"strings"
)

func ToAuthorizationRedis(authorization dto.Authorization) models.Authorization {
	return models.Authorization{
		Code:          authorization.Code(),
		Username:      authorization.Username(),
		ClientID:      authorization.ClientID(),
		RedirectURI:   authorization.RedirectURI(),
		CodeChallenge: authorization.CodeChallenge(),
		Scope:         strings.Join(authorization.Scope(), " "),
	}
}

func ToAuthorizationDTO(authorization models.Authorization) dto.Authorization {
	return dto.NewAuthorization(
		authorization.Code,
		authorization.Username,
		authorization.ClientID,
		authorization.RedirectURI,
		authorization.CodeChallenge,
		strings.Split(authorization.Scope, " "),
	)
}
