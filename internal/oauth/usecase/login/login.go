package login

import (
	"context"
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"log/slog"
	"time"
)

// Repository is a repository for log in use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string) (dto.Client, error)
	UserByUsername(ctx context.Context, username string) (dto.User, error)
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	SaveFlow(ctx context.Context, flow dto.Flow, ttl time.Duration) error
}

// UseCase is a use-case for logging in a user.
type UseCase struct {
	log   *slog.Logger
	repo  Repository
	redis Redis
}

// New returns new log in use-case.
func New(log *slog.Logger, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}

// Execute executes the use-case for logging in a user.
func (uc *UseCase) Execute(ctx context.Context, data Params) (string, *domain.DisplayableError) {
	const op = "usecase.auth.login"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("flow_id", data.FlowID),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		slog.String("scope", data.Scope),
		slog.String("username", data.Username),
		slog.String("password", data.Password),
	)
	log.Info("attempting to login user")

	uriBuilder := builder.NewURI()

	// get flow from redis
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return "", domain.InternalError(err)
		}
	}

	// get client from storage
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return "", domain.InternalError(err)
		}
	}

	// get user from storage
	user, err := uc.repo.UserByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return "", domain.InternalError(err)
		}
	}

	// login logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uriBuilder),
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
	)
	if displayableErr := oauthEntity.Login(loginParams); displayableErr != nil {
		if displayableErr.IsInternal() {
			log.Error(displayableErr.Error())
		}

		return "", displayableErr
	}

	// update flow data in redis
	flow, err = oauthEntity.Flow()
	if err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	if err = uc.redis.SaveFlow(ctx, flow, oauth.RedisAuthorizationCodeTTL*time.Minute); err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	log.Info("user logged in successfully")

	return oauthEntity.RedirectURI(), nil
}
