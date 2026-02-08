package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

// Repository is a repository for exchange token use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string) (dto.Client, error)
	UserByUsername(ctx context.Context, username string) (dto.User, error)
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	RemoveFlow(ctx context.Context, id uuid.UUID) error
}

// UseCase is a use-case for exchange token.
type UseCase struct {
	log   *slog.Logger
	repo  Repository
	redis Redis
}

// New returns new exchange token use-case.
func New(log *slog.Logger, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}

// Execute executes the use-case for exchange token.
func (uc *UseCase) Execute(ctx context.Context, data Params) (dto.Token, error) {
	const op = "usecase.auth.token"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("flow_id", data.FlowID),
		slog.String("grant_type", data.GrantType),
		slog.String("client_id", data.ClientID),
		slog.String("authorization_code", data.AuthorizationCode),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("code_verifier", data.CodeVerifier),
		sl.Strings("scope", data.Scope),
	)
	log.Info("attempting to exchange token")

	// get flow from redis.
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return dto.Token{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	// get user from storage.
	user, err := uc.repo.UserByUsername(ctx, flow.Username())
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return dto.Token{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	// get client from storage.
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn(err.Error())
		} else {
			log.Error(err.Error())

			return dto.Token{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	// exchange token logic.
	oauthEntity := entity.NewOAuth(
		entity.WithFlow(flow),
		entity.WithUser(user),
		entity.WithClient(client),
	)

	exchangeTokenParams := dto.NewExchangeToken(
		data.FlowID,
		data.GrantType,
		data.ClientID,
		data.AuthorizationCode,
		data.RedirectURI,
		data.CodeVerifier,
		data.Scope,
	)
	tokens, err := oauthEntity.ExchangeToken(exchangeTokenParams)
	if err != nil {
		log.Error(err.Error())

		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	// remove flow from redis.
	if err = uc.redis.RemoveFlow(ctx, flow.ID()); err != nil {
		log.Error(err.Error())

		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}
