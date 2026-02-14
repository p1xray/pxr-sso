package authorize

import (
	"context"
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"log/slog"
)

type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

type Redis interface {
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for OAuth authorize.
type UseCase struct {
	log   *slog.Logger
	repo  Repository
	redis Redis
}

// New returns new OAuth authorize use-case.
func New(log *slog.Logger, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}

// Execute executes the use-case for OAuth authorize.
func (uc *UseCase) Execute(ctx context.Context, data Params) string {
	const op = "usecase.oauth.authorize"

	log := uc.log.With(
		slog.String("op", op),
		sl.Strings("response_type", data.ResponseType),
		sl.Strings("client_id", data.ClientID),
		sl.Strings("redirect_uri", data.RedirectURI),
		sl.Strings("code_challenge", data.CodeChallenge),
		sl.Strings("code_challenge_method", data.CodeChallengeMethod),
		sl.Strings("state", data.State),
		sl.Strings("scope", data.Scope),
	)
	log.Info("attempting to initiate user authorization")

	uriBuilder := builder.NewURI()

	// get client from storage
	nullableClient := nullable.None[dto.Client]()
	if len(data.ClientID) == 1 {
		client, err := uc.repo.ClientByCode(ctx, data.ClientID[0])
		if err != nil {
			if errors.Is(err, infrastructure.ErrEntityNotFound) {
				log.Warn("client not found by code %s", data.ClientID[0])
			} else {
				log.Error(err.Error())

				redirectURI := uriBuilder.BuildErrorRedirectURI("", domain.ServerErrorOAuthError(err))

				return redirectURI
			}
		}

		nullableClient = nullable.Some(client)
	}

	// authorize logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uriBuilder),
		entity.WithNullableClient(nullableClient),
	)

	authorizeParams := dto.NewAuthorize(
		data.ResponseType,
		data.ClientID,
		data.RedirectURI,
		data.CodeChallenge,
		data.CodeChallengeMethod,
		data.State,
		data.Scope,
	)
	err := oauthEntity.Authorize(authorizeParams)
	if err != nil {
		log.Warn("%s: %s", "initiate user authorization", err.Error())

		return oauthEntity.RedirectURI()
	}

	// save flow data to redis
	flow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(err.Error())

		oauthEntity.HandleError(domain.ServerErrorOAuthError(err), "")
		return oauthEntity.RedirectURI()
	}

	if err = uc.redis.SaveFlow(ctx, flow); err != nil {
		log.Error(err.Error())

		oauthEntity.HandleError(domain.ServerErrorOAuthError(err), "")
		return oauthEntity.RedirectURI()
	}

	log.Info("initiate user authorization successfully")

	return oauthEntity.RedirectURI()
}
