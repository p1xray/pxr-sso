package converter

import (
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

func ToSessionGrantedScopeLinkStorage(
	sessionID int64,
	scopeID int64,
	setters ...models.SessionGrantedScopeLinkOption,
) models.SessionGrantedScopeLink {
	link := models.SessionGrantedScopeLink{
		SessionID: sessionID,
		ScopeID:   scopeID,
	}

	for _, setter := range setters {
		setter(&link)
	}

	return link
}
