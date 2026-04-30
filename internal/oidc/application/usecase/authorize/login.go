package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type LoginURIBuilder interface {
	BuildLoginRedirectURI(requestURI string) (string, error)
}

type LoginAuthorizationRequestSaver interface {
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

type loginFlowProcessor struct {
	uriBuilder                LoginURIBuilder
	authorizationRequestSaver LoginAuthorizationRequestSaver
}

// NewLoginFlowProcessor returns new authorize login use-case.
func NewLoginFlowProcessor(
	uriBuilder LoginURIBuilder,
	authorizationRequestSaver LoginAuthorizationRequestSaver,
) *loginFlowProcessor {
	return &loginFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithLogin executes the authorize with login.
func (l *loginFlowProcessor) AuthorizeWithLogin(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := l.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return l.uriBuilder.BuildLoginRedirectURI(requestURI)
}
