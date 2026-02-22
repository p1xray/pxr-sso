package register

import (
	"context"
	"errors"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

const logTag = "[pxr-sso-register-use-case]"

// Repository is a repository for register a new user use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
	UserByUsername(ctx context.Context, username string, opts ...repository.UserOption) (dto.User, error)
	CreateUser(ctx context.Context, user dto.User, clientID int64) error
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	SaveFlow(ctx context.Context, flow dto.Flow) error
}

// UseCase is a use-case for registering a new user.
type UseCase struct {
	log        *slog.Logger
	uriBuilder *builder.URI
	repo       Repository
	redis      Redis
}

// New returns new register a new user use-case.
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
	log := uc.log.With(
		slog.String("flow_id", data.FlowID),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		sl.Strings("scope", data.Scope),
		slog.String("username", data.Username),
		slog.String("full_name", data.FullName),
	)
	log.Debug(logTag + " attempting to register new user")

	// get flow from redis
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get flow", sl.Err(err))

		return "", err
	}

	// get client from storage
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return "", err
	}

	// get user from storage
	user, err := uc.repo.UserByUsername(ctx, data.Username)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get user by username", sl.Err(err))

		return "", err
	}

	// login logic
	oauthEntity := entity.NewOAuth(
		entity.WithBuilderURI(uc.uriBuilder),
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
	if err = oauthEntity.Register(registerParams); err != nil {
		var validationErr *validator.Error
		if errors.As(err, &validationErr) && !validationErr.IsInvalidUserCredentials() {
			log.Error(logTag+" register process", sl.Err(err))
		}

		return "", err
	}

	// create new user in storage
	newUser, err := oauthEntity.User()
	if err != nil {
		log.Error(logTag+" get the generated user", sl.Err(err))

		return "", err
	}

	if err = uc.repo.CreateUser(ctx, newUser, client.ID()); err != nil {
		log.Error(logTag+" save user", sl.Err(err))

		return "", err
	}

	// update flow data in redis
	updatedFlow, err := oauthEntity.Flow()
	if err != nil {
		log.Error(logTag+" get the updated flow", sl.Err(err))

		return "", err
	}

	if err = uc.redis.SaveFlow(ctx, updatedFlow); err != nil {
		log.Error(logTag+" save flow", sl.Err(err))

		return "", err
	}

	log.Debug(logTag + " user register successfully")

	return oauthEntity.RedirectURI(), nil
}
