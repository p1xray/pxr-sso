package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

// ConsentURIBuilder is the URI builder to redirect the user agent to the confirming consent page.
type ConsentURIBuilder interface {
	// BuildConsentRedirectURI return the URI to redirect to the confirming consent page.
	BuildConsentRedirectURI(requestURI string) (string, error)
}

// ConsentAuthorizationRequestSaver is the saver authorization request to a storage.
type ConsentAuthorizationRequestSaver interface {
	// SaveAuthorizationRequest saves the authorization request to a storage.
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

// consentFlowProcessor is the processor for authorize flow with confirming consent interaction.
type consentFlowProcessor struct {
	uriBuilder                ConsentURIBuilder
	authorizationRequestSaver ConsentAuthorizationRequestSaver
}

// NewConsentFlowProcessor creates a new processor for authorize flow with confirming consent interaction.
func NewConsentFlowProcessor(
	uriBuilder ConsentURIBuilder,
	authorizationRequestSaver ConsentAuthorizationRequestSaver,
) *consentFlowProcessor {
	return &consentFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithConsent processes the authorize flow with confirming consent interaction.
//
// This method saves the authorization request to a storage
// and generates the uri to redirect to the confirming consent page.
func (c *consentFlowProcessor) AuthorizeWithConsent(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := c.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return c.uriBuilder.BuildConsentRedirectURI(requestURI)
}
