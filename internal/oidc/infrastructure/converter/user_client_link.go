package converter

import "github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"

func ToUserClientLinkStorage(userID, clientID int64, setters ...models.UserClientLinkOption) models.UserClientLink {
	userClientLinkModel := models.UserClientLink{
		UserID:   userID,
		ClientID: clientID,
	}

	for _, setter := range setters {
		setter(&userClientLinkModel)
	}

	return userClientLinkModel
}
