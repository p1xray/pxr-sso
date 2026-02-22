package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

const (
	logTag = "[pxr-sso-exchange-token-use-case]"
	pkgTag = "token use case"
)

// Repository is a repository for exchange token use-case.
type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
	UserByUsername(ctx context.Context, username string, opts ...repository.UserOption) (dto.User, error)
}

type Redis interface {
	Flow(ctx context.Context, id string) (dto.Flow, error)
	RemoveFlow(ctx context.Context, id uuid.UUID) error
}

// UseCase is a use-case for exchange token.
type UseCase struct {
	log            *slog.Logger
	tokenGenerator *generator.Token
	repo           Repository
	redis          Redis
}

// New returns new exchange token use-case.
func New(log *slog.Logger, tokenGenerator *generator.Token, repo Repository, redis Redis) *UseCase {
	return &UseCase{
		log:            log,
		tokenGenerator: tokenGenerator,
		repo:           repo,
		redis:          redis,
	}
}

// Execute executes the use-case for exchange token.
func (uc *UseCase) Execute(ctx context.Context, data Params) (oauth.ServiceData[dto.Token], error) {
	output, err := oauth.Call(func() (dto.Token, error) {
		return uc.exchangeToken(ctx, data)
	})

	return output, err
}

func (uc *UseCase) exchangeToken(ctx context.Context, data Params) (dto.Token, error) {
	const op = "exchange token"

	log := uc.log.With(
		slog.String("flow_id", data.FlowID),
		slog.String("grant_type", data.GrantType),
		slog.String("client_id", data.ClientID),
		slog.String("authorization_code", data.AuthorizationCode),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("code_verifier", data.CodeVerifier),
		sl.Strings("scope", data.Scope),
	)
	log.Debug(logTag + " attempting to exchange token")

	// get flow from redis.
	flow, err := uc.redis.Flow(ctx, data.FlowID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get flow", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// get user from storage.
	user, err := uc.repo.UserByUsername(ctx, flow.Username())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get user by username", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// get client from storage.
	client, err := uc.repo.ClientByCode(ctx, data.ClientID)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// exchange token logic.
	oauthEntity := entity.NewOAuth(
		entity.WithTokenGenerator(uc.tokenGenerator),
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
		data.Audience,
		data.Scope,
	)
	tokens, err := oauthEntity.ExchangeToken(exchangeTokenParams)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// remove flow from redis.
	if err = uc.redis.RemoveFlow(ctx, flow.ID()); err != nil {
		log.Error(logTag+" save flow", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return tokens, nil
}
