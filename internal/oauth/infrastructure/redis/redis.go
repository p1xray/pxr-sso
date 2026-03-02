package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/models"
	"github.com/redis/go-redis/v9"
	"time"
)

const pkgTag = "redis storage"

type Redis struct {
	client *redis.Client

	flowTTL          time.Duration
	authorizationTTL time.Duration
}

func New(cfg Config) (*Redis, error) {
	opt, err := redis.ParseURL(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", pkgTag, "parse connection URL", err)
	}
	client := redis.NewClient(opt)

	return &Redis{
		client:           client,
		flowTTL:          cfg.FlowTTL,
		authorizationTTL: cfg.AuthorizationTTL,
	}, nil
}

func (r *Redis) Flow(ctx context.Context, id string) (dto.Flow, error) {
	const op = "get flow"

	redisFlowKey := builder.BuildRedisFlowKey(id)
	cmd := r.client.Get(ctx, redisFlowKey)

	redisFlow := models.Flow{}
	if err := cmd.Scan(&redisFlow); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.Flow{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return dto.Flow{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	flow, err := converter.ToFlowDTO(redisFlow)
	if err != nil {
		return dto.Flow{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return flow, nil
}

func (r *Redis) SaveFlow(ctx context.Context, flow dto.Flow) error {
	const op = "save flow"

	redisFlow := converter.ToFlowRedis(flow)
	redisFlowKey := builder.BuildRedisFlowKey(redisFlow.ID)

	cmd := r.client.Set(ctx, redisFlowKey, redisFlow, r.flowTTL)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (r *Redis) RemoveFlow(ctx context.Context, id uuid.UUID) error {
	const op = "remove flow"

	redisFlowKey := builder.BuildRedisFlowKey(id.String())

	cmd := r.client.Del(ctx, redisFlowKey)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (r *Redis) Authorization(ctx context.Context, id string) (dto.Authorization, error) {
	const op = "get authorization"

	cmd := r.client.Get(ctx, id)

	redisAuthorization := models.Authorization{}
	if err := cmd.Scan(&redisAuthorization); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.Authorization{}, fmt.Errorf("%s: %s: %w", pkgTag, op, infrastructure.ErrEntityNotFound)
		}

		return dto.Authorization{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	authorization := converter.ToAuthorizationDTO(redisAuthorization)
	return authorization, nil
}

func (r *Redis) SaveAuthorization(ctx context.Context, authorization dto.Authorization) error {
	const op = "save authorization"

	redisAuthorization := converter.ToAuthorizationRedis(authorization)

	cmd := r.client.Set(ctx, redisAuthorization.Code, redisAuthorization, r.authorizationTTL)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}

func (r *Redis) RemoveAuthorization(ctx context.Context, id string) error {
	const op = "remove authorization"

	cmd := r.client.Del(ctx, id)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return nil
}
