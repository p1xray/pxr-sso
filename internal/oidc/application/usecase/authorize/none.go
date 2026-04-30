package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type AuthorizedGrantProcessor interface {
	Execute(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type noneFlowProcessor struct {
	authorizedGrantProcessor AuthorizedGrantProcessor
}

// NewNoneFlowProcessor returns new authorize none use-case.
func NewNoneFlowProcessor(
	authorizedGrantProcessor AuthorizedGrantProcessor,
) *noneFlowProcessor {
	return &noneFlowProcessor{
		authorizedGrantProcessor: authorizedGrantProcessor,
	}
}

// AuthorizeWithoutInteraction executes the authorize without interaction.
func (n *noneFlowProcessor) AuthorizeWithoutInteraction(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	return n.authorizedGrantProcessor.Execute(ctx, data)
}
