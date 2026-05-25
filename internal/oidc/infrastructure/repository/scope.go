package repository

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/converter"
)

func (r *repository) ScopesByCode(ctx context.Context, codes []string) ([]dto.Scope, error) {
	const op = "get scopes by code"

	scopes, err := r.storage.ScopesByCode(ctx, codes)
	if err != nil {
		return []dto.Scope{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	scopesDTO := make([]dto.Scope, len(scopes))
	for i, scope := range scopes {
		scopesDTO[i] = converter.ToScopeDTO(scope)
	}

	return scopesDTO, nil
}
