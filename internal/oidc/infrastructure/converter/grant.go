package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache/models"
)

func ToAuthorizedGrantDTO(grant models.AuthorizedGrant) dto.AuthorizedGrant {
	return dto.NewAuthorizedGrant(
		ToAuthorizedSessionDTO(grant.Session),
		ToAuthorizeRequestDTO(grant.Request),
	)
}

func ToAuthorizedSessionDTO(session models.Session) dto.AuthorizedSession {
	return dto.NewAuthorizedSession(
		session.ID,
		session.Subject,
		session.AuthTime,
		session.IdentityProvider,
	)
}

func ToAuthorizedGrantStorage(grant dto.AuthorizedGrant) models.AuthorizedGrant {
	session := grant.Session()
	request := grant.Request()

	return models.AuthorizedGrant{
		Session: ToAuthorizedSessionStorage(session),
		Request: ToAuthorizationRequestStorage(request),
	}
}

func ToAuthorizedSessionStorage(session dto.AuthorizedSession) models.Session {
	return models.Session{
		ID:               session.ID(),
		Subject:          session.Subject(),
		AuthTime:         session.AuthTime(),
		IdentityProvider: session.IdentityProvider(),
	}
}
