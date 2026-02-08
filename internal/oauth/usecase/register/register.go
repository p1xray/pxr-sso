package register

import (
	"context"
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

// Repository is a repository for register a new user use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string) (dto.Client, error)
	UserByUsername(ctx context.Context, username string) (dto.User, error)
	CreateUser(ctx context.Context, user dto.User, clientID int64) error
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for registering a new user.
type UseCase struct {
	log   *slog.Logger
	repo  Repository
	redis Redis
}

// New returns new register a new user use-case.
func New(log *slog.Logger, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}

// Execute executes the use-case for registering a new user.
func (uc *UseCase) Execute(ctx context.Context, data Params) (string, *domain.DisplayableError) {
	const op = "usecase.auth.register"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("flow_id", data.FlowID),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		sl.Strings("scope", data.Scope),
		slog.String("username", data.Username),
		slog.String("full_name", data.FullName),
	)
	log.Info("attempting to register new user")

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

	registerParams := dto.NewRegister(
		data.FlowID,
		data.ResponseType,
		data.ClientID,
		data.RedirectURI,
		data.State,
		data.Username,
		data.Password,
		data.FullName,
		data.Scope,
	)
	if displayableErr := oauthEntity.Register(registerParams); displayableErr != nil {
		if displayableErr.IsInternal() {
			log.Error(displayableErr.Error())
		}

		return "", displayableErr
	}

	// create new user in storage
	newUser, err := oauthEntity.User()
	if err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	if err = uc.repo.CreateUser(ctx, newUser, client.ID); err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	// update flow data in redis
	updatedFlow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	if err = uc.redis.SaveFlow(ctx, updatedFlow); err != nil {
		log.Error(err.Error())

		return "", domain.InternalError(err)
	}

	log.Info("user register successfully")

	return oauthEntity.RedirectURI(), nil
}
