package login

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
	"log/slog"
)

const logTag = "[pxr-sso-login-use-case]"

type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
	UserByUsername(ctx context.Context, username string, opts ...repository.UserOption) (dto.User, error)
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for logging in a user.
type UseCase struct {
	log        *slog.Logger
	uriBuilder *builder.URI
	repo       Repository
	redis      Redis
}

// New returns new log in use-case.
func New(log *slog.Logger, uriBuilder *builder.URI, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:        log,
		uriBuilder: uriBuilder,
		repo:       repo,
		redis:      redis,
	}
}

// Execute executes the use-case for logging in a user.
func (uc *UseCase) Execute(ctx context.Context, data Params) (string, *domain.DisplayableError) {
	log := uc.log.With(
		slog.String("flow_id", data.FlowID),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		sl.Strings("scope", data.Scope),
		slog.String("username", data.Username),
	)
	log.Debug(logTag + " attempting to login user")

	// get flow from redis
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get flow", sl.Err(err))

		return "", domain.InternalError(err)
	}

	// get client from storage
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return "", domain.InternalError(err)
	}

	// get user from storage
	user, err := uc.repo.UserByUsername(ctx, data.Username)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get user by username", sl.Err(err))

		return "", domain.InternalError(err)
	}

	// login logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uc.uriBuilder),
		entity.WithFlow(flow),
		entity.WithClient(client),
		entity.WithUser(user),
	)

	loginParams := dto.NewLogin(
		data.FlowID,
		data.ResponseType,
		data.ClientID,
		data.RedirectURI,
		data.State,
		data.Username,
		data.Password,
		data.Scope,
	)
	if displayableErr := oauthEntity.Login(loginParams); displayableErr != nil {
		if displayableErr.IsInternal() {
			log.Error(logTag+" login process", sl.Err(displayableErr.Unwrap()))
		}

		return "", displayableErr
	}

	// update flow data in redis
	updatedFlow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(logTag+" get the updated flow", sl.Err(err))

		return "", domain.InternalError(err)
	}

	if err = uc.redis.SaveFlow(ctx, updatedFlow); err != nil {
		log.Error(logTag+" save flow", sl.Err(err))

		return "", domain.InternalError(err)
	}

	log.Debug(logTag + " user logged in successfully")

	return oauthEntity.RedirectURI(), nil
}
