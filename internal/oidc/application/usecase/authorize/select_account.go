package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type SelectAccountURIBuilder interface {
	BuildSelectAccountRedirectURI(requestURI string) (string, error)
}

type SelectAccountAuthorizationRequestSaver interface {
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

type selectAccountFlowProcessor struct {
	uriBuilder                SelectAccountURIBuilder
	authorizationRequestSaver SelectAccountAuthorizationRequestSaver
}

// NewSelectAccountFlowProcessor returns new authorize select account use-case.
func NewSelectAccountFlowProcessor(
	uriBuilder SelectAccountURIBuilder,
	authorizationRequestSaver SelectAccountAuthorizationRequestSaver,
) *selectAccountFlowProcessor {
	return &selectAccountFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithSelectAccount executes the authorize with select account.
func (sa *selectAccountFlowProcessor) AuthorizeWithSelectAccount(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := sa.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return sa.uriBuilder.BuildSelectAccountRedirectURI(requestURI)
}
