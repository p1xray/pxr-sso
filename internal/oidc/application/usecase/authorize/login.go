package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

// LoginURIBuilder is the URI builder to redirect the user agent to the login page.
type LoginURIBuilder interface {
	// BuildLoginRedirectURI return the URI to redirect to the login page.
	BuildLoginRedirectURI(requestURI string) (string, error)
}

// LoginAuthorizationRequestSaver is the saver authorization request to a storage.
type LoginAuthorizationRequestSaver interface {
	// SaveAuthorizationRequest saves the authorization request to a storage.
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

// loginFlowProcessor is the processor for authorize flow with logging in interaction.
type loginFlowProcessor struct {
	uriBuilder                LoginURIBuilder
	authorizationRequestSaver LoginAuthorizationRequestSaver
}

// NewLoginFlowProcessor creates a new processor for authorize flow with logging in interaction.
func NewLoginFlowProcessor(
	uriBuilder LoginURIBuilder,
	authorizationRequestSaver LoginAuthorizationRequestSaver,
) *loginFlowProcessor {
	return &loginFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithLogin processes the authorize flow with logging in interaction.
//
// This method saves the authorization request to a storage
// and generates the uri to redirect to the login page.
func (l *loginFlowProcessor) AuthorizeWithLogin(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := l.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return l.uriBuilder.BuildLoginRedirectURI(requestURI)
}
