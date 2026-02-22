package authorize

import (
	"context"
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"log/slog"
)

const logTag = "[pxr-sso-authorize-use-case]"

type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

type Redis interface {
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for OAuth authorize.
type UseCase struct {
	log        *slog.Logger
	uriBuilder *builder.URI
	repo       Repository
	redis      Redis
}

// New returns new OAuth authorize use-case.
func New(log *slog.Logger, uriBuilder *builder.URI, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:        log,
		uriBuilder: uriBuilder,
		repo:       repo,
		redis:      redis,
	}
}

// Execute executes the use-case for OAuth authorize.
func (uc *UseCase) Execute(ctx context.Context, data Params) (oauth.ServiceData[string], error) {
	output, err := oauth.Call(func() (string, error) {
		return uc.authorize(ctx, data)
	})

	return output, err
}

func (uc *UseCase) authorize(ctx context.Context, data Params) (string, error) {
	log := uc.log.With(
		sl.Strings("response_type", data.ResponseType),
		sl.Strings("client_id", data.ClientID),
		sl.Strings("redirect_uri", data.RedirectURI),
		sl.Strings("code_challenge", data.CodeChallenge),
		sl.Strings("code_challenge_method", data.CodeChallengeMethod),
		sl.Strings("state", data.State),
		sl.Strings("scope", data.Scope),
	)
	log.Debug(logTag + " attempting to initiate user authorization")

	// get client from storage
	nullableClient := nullable.None[dto.Client]()
	if len(data.ClientID) == 1 {
		client, err := uc.repo.ClientByCode(ctx, data.ClientID[0])
		if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Error(logTag+" get client by code", sl.Err(err))

			redirectURI := uc.uriBuilder.BuildErrorRedirectURI("", err)
			return redirectURI, err

		}

		nullableClient = nullable.Some(client)
	}

	// authorize logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uc.uriBuilder),
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
		return oauthEntity.RedirectURI(), err
	}

	// save flow data to redis
	flow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(logTag+" get the generated flow", sl.Err(err))

		oauthEntity.HandleError(err, "")
		return oauthEntity.RedirectURI(), err
	}

	if err = uc.redis.SaveFlow(ctx, flow); err != nil {
		log.Error(logTag+" save flow", sl.Err(err))

		oauthEntity.HandleError(err, "")
		return oauthEntity.RedirectURI(), err
	}

	log.Debug(logTag + " initiate user authorization successfully")

	return oauthEntity.RedirectURI(), err
}
