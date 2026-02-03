package login

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"log/slog"
	"time"
)

// Repository is a repository for log in use-case.
type Repository interface {
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
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	// TODO: get client from storage
	nullableClient := nullable.None[dto.Client]()

	// TODO: get user from storage
	nullableUser := nullable.None[dto.User]()

	// login logic
	oauthEntity := entity.NewOAuth(
		uriBuilder,
		entity.WithFlow(flow),
		entity.WithClient(nullableClient),
		entity.WithUser(nullableUser),
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
	if oauthErr := oauthEntity.Login(loginParams); oauthErr != nil {
		if oauthErr.IsInternal() {
			log.Error(oauthErr.Error())
		}

		return "", oauthErr
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
