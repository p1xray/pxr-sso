package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache/builder"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache/models"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/converter"
	"github.com/redis/go-redis/v9"
	"time"
)

const pkgTag = "redis storage"

type cache struct {
	client *redis.Client

	authorizationRequestTTL time.Duration
	authorizedGrantTTL      time.Duration
}

func New(cfg Config) (*cache, error) {
	opt, err := redis.ParseURL(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, "parse connection URL", err)
	}
	client := redis.NewClient(opt)

	return &cache{
		client:                  client,
		authorizationRequestTTL: cfg.AuthorizationRequestTTL,
		authorizedGrantTTL:      cfg.AuthorizedGrantTTL,
	}, nil
}

func (c *cache) AuthorizationRequest(
	ctx context.Context,
	requestURI string,
) (dto.ValidatedAuthorizeRequest, error) {
	const op = "get authorization request from redis storage"

	key := builder.AuthorizationRequestKey(requestURI)
	cmd := c.client.Get(ctx, key)

	redisAuthorizationRequest := models.AuthorizationRequest{}
	if err := cmd.Scan(&redisAuthorizationRequest); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.ValidatedAuthorizeRequest{},
				fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return dto.ValidatedAuthorizeRequest{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	authorizationRequest := converter.ToAuthorizeRequestDTO(redisAuthorizationRequest)

	return authorizationRequest, nil
}

func (c *cache) SaveAuthorizationRequest(
	ctx context.Context,
	requestURI string,
	request dto.ValidatedAuthorizeRequest,
) error {
	const op = "save authorization request to redis storage"

	key := builder.AuthorizationRequestKey(requestURI)
	redisAuthorizationRequest := converter.ToAuthorizationRequestStorage(request)

	cmd := c.client.Set(ctx, key, redisAuthorizationRequest, c.authorizationRequestTTL)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (c *cache) RemoveAuthorizationRequest(
	ctx context.Context,
	requestURI string,
) error {
	const op = "remove authorization request from redis storage"

	key := builder.AuthorizationRequestKey(requestURI)
	cmd := c.client.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (c *cache) AuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
) (dto.AuthorizedGrant, error) {
	const op = "get authorized grant from redis storage"

	key := builder.AuthorizedGrantKey(authorizationCode)
	cmd := c.client.Get(ctx, key)

	redisAuthorizedGrant := models.AuthorizedGrant{}
	if err := cmd.Scan(&redisAuthorizedGrant); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.AuthorizedGrant{},
				fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return dto.AuthorizedGrant{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	authorizedGrant := converter.ToAuthorizedGrantDTO(redisAuthorizedGrant)

	return authorizedGrant, nil
}

func (c *cache) SaveAuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
	authorizedGrant dto.AuthorizedGrant,
) error {
	const op = "save authorized grant to redis storage"

	key := builder.AuthorizedGrantKey(authorizationCode)
	redisAuthorizedGrant := converter.ToAuthorizedGrantStorage(authorizedGrant)

	cmd := c.client.Set(ctx, key, redisAuthorizedGrant, c.authorizedGrantTTL)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (c *cache) RemoveAuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
) error {
	const op = "remove authorized grant from redis storage"

	key := builder.AuthorizedGrantKey(authorizationCode)
	cmd := c.client.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}
