package token

import (
	"context"
	"errors"
	"fmt"
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
	Authorization(ctx context.Context, id string) (dto.Authorization, error)
	RemoveAuthorization(ctx context.Context, id string) error
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
		slog.String("grant_type", data.GrantType),
		slog.String("client_id", data.ClientID),
		slog.String("authorization_code", data.AuthorizationCode),
		slog.String("redirect_uri", data.RedirectURI),
		slog.String("code_verifier", data.CodeVerifier),
		slog.String("audience", data.Audience),
		sl.Strings("scope", data.Scope),
	)
	log.Debug(logTag + " attempting to exchange token")

	// get authorization from redis.
	authorization, err := uc.redis.Authorization(ctx, data.AuthorizationCode)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get authorization", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// get user from storage.
	user, err := uc.repo.UserByUsername(ctx, authorization.Username(), repository.WithRoles())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get user by username", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// get client from storage.
	client, err := uc.repo.ClientByCode(
		ctx,
		data.ClientID,
		repository.WithAudiences(),
		repository.WithRedirectURIs(),
		repository.WithScopes())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	// exchange token logic.
	oauthEntity := entity.NewOAuth(
		entity.WithTokenGenerator(uc.tokenGenerator),
		entity.WithAuthorization(authorization),
		entity.WithUser(user),
		entity.WithClient(client),
	)

	exchangeTokenParams := dto.NewExchangeToken(
		data.GrantType,
		data.ClientID,
		data.ClientSecret,
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

	// remove authorization data from redis.
	authorizationToRemove, err := oauthEntity.Authorization()
	if err != nil {
		log.Error(logTag+" get authorization", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	if err = uc.redis.RemoveAuthorization(ctx, authorizationToRemove.Code()); err != nil {
		log.Error(logTag+" remove authorization", sl.Err(err))

		return dto.Token{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	log.Debug(logTag + " exchanging tokens successfully")

	return tokens, nil
}
