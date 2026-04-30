package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type ConsentURIBuilder interface {
	BuildConsentRedirectURI(requestURI string) (string, error)
}

type ConsentAuthorizationRequestSaver interface {
	SaveAuthorizationRequest(ctx context.Context, requestURI string, request dto.ValidatedAuthorizeRequest) error
}

type consentFlowProcessor struct {
	uriBuilder                ConsentURIBuilder
	authorizationRequestSaver ConsentAuthorizationRequestSaver
}

// NewConsentFlowProcessor returns new authorize consent use-case.
func NewConsentFlowProcessor(
	uriBuilder ConsentURIBuilder,
	authorizationRequestSaver ConsentAuthorizationRequestSaver,
) *consentFlowProcessor {
	return &consentFlowProcessor{
		uriBuilder:                uriBuilder,
		authorizationRequestSaver: authorizationRequestSaver,
	}
}

// AuthorizeWithConsent executes the authorize with consent.
func (c *consentFlowProcessor) AuthorizeWithConsent(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	requestURI := generator.RequestURI()
	if err := c.authorizationRequestSaver.SaveAuthorizationRequest(ctx, requestURI, data.Request()); err != nil {
		return "", err
	}

	return c.uriBuilder.BuildConsentRedirectURI(requestURI)
}
