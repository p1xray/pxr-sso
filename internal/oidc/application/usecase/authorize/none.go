package authorize

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

// AuthorizedGrantProcessor is the processor for processing authorized grant.
type AuthorizedGrantProcessor interface {
	// Execute processing authorized grant.
	Execute(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// noneFlowProcessor is the processor for authorize without interaction flow.
type noneFlowProcessor struct {
	authorizedGrantProcessor AuthorizedGrantProcessor
}

// NewNoneFlowProcessor creates a new processor for authorize without interaction flow.
func NewNoneFlowProcessor(
	authorizedGrantProcessor AuthorizedGrantProcessor,
) *noneFlowProcessor {
	return &noneFlowProcessor{
		authorizedGrantProcessor: authorizedGrantProcessor,
	}
}

// AuthorizeWithoutInteraction processes the authorize without interaction flow.
//
// This method just processing authorized grant.
func (n *noneFlowProcessor) AuthorizeWithoutInteraction(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	return n.authorizedGrantProcessor.Execute(ctx, data)
}
