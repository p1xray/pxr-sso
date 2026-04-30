package cache

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type Cache interface {
	AuthorizationRequest(ctx context.Context, requestURI string) (dto.ValidatedAuthorizeRequest, error)
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
	RemoveAuthorizationRequest(ctx context.Context, requestURI string) error

	AuthorizedGrant(ctx context.Context, authorizationCode string) (dto.AuthorizedGrant, error)
	SaveAuthorizedGrant(ctx context.Context, authorizationCode string, authorizedGrant dto.AuthorizedGrant) error
	RemoveAuthorizedGrant(ctx context.Context, authorizationCode string) error
}
