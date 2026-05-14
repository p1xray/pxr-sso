package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToSessionsDTO(sessions []models.Session) []dto.Session {
	sessionsDTO := make([]dto.Session, len(sessions))

	for i, session := range sessions {
		sessionsDTO[i] = ToSessionDTO(session)
	}

	return sessionsDTO
}

func ToSessionDTO(session models.Session) dto.Session {
	scopes := make([]dto.Scope, len(session.GrantedScopeLinks))
	for i, link := range session.GrantedScopeLinks {
		scopes[i] = ToScopeDTO(link.Scope)
	}

	return dto.NewSession(
		session.ID,
		session.Code,
		ToClientDTO(session.Client),
		ToUserDTO(session.User),
		session.AuthTime,
		session.IdentityProvider,
		scopes,
	)
}

func ToSessionStorage(session entity.Session, setters ...models.SessionOption) models.Session {
	client := session.Client()
	user := session.User()

	sessionStorageModel := models.Session{
		ID:               session.ID(),
		Code:             session.Code(),
		ClientID:         client.ID(),
		UserID:           user.ID(),
		AuthTime:         session.AuthTime(),
		IdentityProvider: session.IdentityProvider(),
	}

	for _, setter := range setters {
		setter(&sessionStorageModel)
	}

	return sessionStorageModel
}
