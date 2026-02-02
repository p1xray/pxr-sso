package login

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"log/slog"
	"time"
)

// Repository is a repository for log in use-case.
type Repository interface {
}

type Redis interface {
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
func (uc *UseCase) Execute(ctx context.Context, data Params) (string, error) {
	const op = "usecase.auth.login"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("response_type", data.ResponseType),
		slog.String("client_id", data.ClientID),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("state", data.State),
		slog.String("scope", data.Scope),
		slog.String("username", data.Username),
		slog.String("password", data.Password),
	)
	log.Info("attempting to login user")

	// TODO: get flow from redis

	// TODO: get client from storage

	// TODO: get user from storage

	// TODO: login logic

	// TODO: update flow data in redis

	log.Info("user logged in successfully")

	return redirectURIWithParams, nil
}
