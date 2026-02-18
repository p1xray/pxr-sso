package consent

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

const (
	logTag = "[pxr-sso-consent-use-case]"
	pkgTag = "consent use case"
)

// Repository is a repository for confirm consent use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for confirming consent.
type UseCase struct {
	log        *slog.Logger
	uriBuilder *builder.URI
	repo       Repository
	redis      Redis
}

// New returns new confirm consent use-case.
func New(log *slog.Logger, uriBuilder *builder.URI, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:        log,
		uriBuilder: uriBuilder,
		repo:       repo,
		redis:      redis,
	}
}

// Execute executes the use-case for registering a new user.
func (uc *UseCase) Execute(ctx context.Context, data Params) (string, error) {
	const op = "confirm consent"

	log := uc.log.With(
		slog.String("flow_id", data.FlowID),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		sl.Strings("scope", data.Scope),
	)
	log.Debug(logTag + " attempting to confirm consent")

	// get flow from redis
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get flow", sl.Err(err))

		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// get client from storage
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// confirm consent logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uc.uriBuilder),
		entity.WithFlow(flow),
		entity.WithClient(client),
	)

	consentParams := dto.NewConsent(
		data.FlowID,
		data.ResponseType,
		data.ClientID,
		data.RedirectURI,
		data.State,
		data.Scope,
	)
	if err = oauthEntity.Consent(consentParams); err != nil {
		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// update flow data in redis
	updatedFlow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(logTag+" get the updated flow", sl.Err(err))

		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	if err = uc.redis.SaveFlow(ctx, updatedFlow); err != nil {
		log.Error(logTag+" save flow", sl.Err(err))

		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	log.Debug(logTag + " confirming consent successfully")

	return oauthEntity.RedirectURI(), nil
}
