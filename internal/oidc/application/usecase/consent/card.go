package consent

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
	"slices"
)

type CardAuthorizationRequestReader interface {
	AuthorizationRequest(ctx context.Context, requestURI string) (dto.ValidatedAuthorizeRequest, error)
}

type CardScopeReader interface {
	ScopesByCode(ctx context.Context, codes []string) ([]dto.Scope, error)
}

type CardReader interface {
	Read(ctx context.Context, request dto.ConsentCardRequest) (dto.ConsentCardResponse, error)
}

type card struct {
	log                        *slog.Logger
	authorizationRequestReader CardAuthorizationRequestReader
	scopeReader                CardScopeReader
}

func NewCardReader(log *slog.Logger,
	authorizationRequestReader CardAuthorizationRequestReader,
	scopeReader CardScopeReader,
) *card {
	return &card{
		log:                        log,
		authorizationRequestReader: authorizationRequestReader,
		scopeReader:                scopeReader,
	}
}

func (c *card) Read(ctx context.Context, request dto.ConsentCardRequest) (dto.ConsentCardResponse, error) {
	const op = "getting data for the consent card"
	const logTag = "[pxr-sso-use-case-consent-card]"

	log := c.log.With(slog.String("request_uri", request.RequestURI()))
	log.Info(logTag + " " + op + " started")

	authorizationRequest, err := c.authorizationRequestReader.AuthorizationRequest(ctx, request.RequestURI())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get authorization request", sl.Err(err))

		return dto.ConsentCardResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	scopeCodes := authorizationRequest.AllScopes()
	scopes, err := c.scopeReader.ScopesByCode(ctx, scopeCodes)
	if err != nil {
		return dto.ConsentCardResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	consentScopes := make([]dto.ConsentScope, 0)
	grantedScopes := authorizationRequest.GrantedScopes()
	for _, scope := range scopes {
		isGranted := slices.Contains(grantedScopes, scope.Code())
		consentScope := dto.NewConsentScope(scope.Code(), scope.Name(), scope.Description(), isGranted)
		consentScopes = append(consentScopes, consentScope)
	}

	consentCard := dto.NewConsentCardResponse(consentScopes)

	log.Info(logTag + " " + op + " successfully finished")
	return consentCard, nil
}
