package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

// SelectAccountURIBuilder is the URI builder to redirect the user agent to the selecting account page.
type SelectAccountURIBuilder interface {
	// BuildSelectAccountRedirectURI return the URI to redirect to the selecting account page.
	BuildSelectAccountRedirectURI(requestURI string) (string, error)
}

// SelectAccountAuthorizationRequestSaver is the saver authorization request to a storage.
type SelectAccountAuthorizationRequestSaver interface {
	// SaveAuthorizationRequest saves the authorization request to a storage.
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

// selectAccountFlowProcessor is the processor for authorize flow with selecting account interaction.
type selectAccountFlowProcessor struct {
	uriBuilder                SelectAccountURIBuilder
	authorizationRequestSaver SelectAccountAuthorizationRequestSaver
}

// NewSelectAccountFlowProcessor creates a new processor for authorize flow with selecting account interaction.
func NewSelectAccountFlowProcessor(
	uriBuilder SelectAccountURIBuilder,
	authorizationRequestSaver SelectAccountAuthorizationRequestSaver,
) *selectAccountFlowProcessor {
	return &selectAccountFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithSelectAccount processes the authorize flow with selecting account interaction.
//
// This method saves the authorization request to a storage
// and generates the uri to redirect to the selecting account page.
func (sa *selectAccountFlowProcessor) AuthorizeWithSelectAccount(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := sa.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return sa.uriBuilder.BuildSelectAccountRedirectURI(requestURI)
}
