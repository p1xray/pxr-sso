package redis

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/redis/go-redis/v9"
	"time"
)

const pkgTag = "redis storage"

type cache struct {
	client *redis.Client

	flowTTL          time.Duration
	authorizationTTL time.Duration
}

func New(cfg Config) (*cache, error) {
	opt, err := redis.ParseURL(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, "parse connection URL", err)
	}
	client := redis.NewClient(opt)

	return &cache{
		client:           client,
		flowTTL:          cfg.FlowTTL,
		authorizationTTL: cfg.AuthorizationTTL,
	}, nil
}

func (c *cache) AuthorizationRequest(
	ctx context.Context,
	requestURI string,
) (dto.ValidatedAuthorizeRequest, error) {
	// TODO: implement this

	//cmd := r.client.Get(ctx, id)
	//
	//redisAuthorization := models.Authorization{}
	//if err := cmd.Scan(&redisAuthorization); err != nil {
	//	if errors.Is(err, redis.Nil) {
	//		return dto.Authorization{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
	//	}
	//
	//	return dto.Authorization{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	//}
	//
	//authorization := converter.ToAuthorizationDTO(redisAuthorization)

	return dto.ValidatedAuthorizeRequest{}, nil
}

func (c *cache) SaveAuthorizationRequest(
	ctx context.Context,
	requestURI string,
	request dto.ValidatedAuthorizeRequest,
) error {
	// TODO: implement this

	//redisAuthorization := converter.ToAuthorizationRedis(authorization)
	//
	//cmd := r.client.Set(ctx, redisAuthorization.Code, redisAuthorization, r.authorizationTTL)
	//if err := cmd.Err(); err != nil {
	//	return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	//}

	return nil
}

func (c *cache) RemoveAuthorizationRequest(
	ctx context.Context,
	requestURI string,
) error {
	// TODO: implement this

	//cmd := r.client.Del(ctx, id)
	//if err := cmd.Err(); err != nil {
	//	return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	//}

	return nil
}

func (c *cache) AuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
) (dto.AuthorizedGrant, error) {
	// TODO: implement this
	return dto.AuthorizedGrant{}, nil
}

func (c *cache) SaveAuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
	authorizedGrant dto.AuthorizedGrant,
) error {
	// TODO: implement this
	return nil
}

func (c *cache) RemoveAuthorizedGrant(
	ctx context.Context,
	authorizationCode string,
) error {
	// TODO: implement this
	return nil
}
