package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
}

func New(connectionURL string) (*Redis, error) {
	opt, err := redis.ParseURL(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "parse redis connection URL", err)
	}
	client := redis.NewClient(opt)

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) Flow(ctx context.Context, id string) (dto.Flow, error) {
	const op = "infrastructure.redis.Flow"

	redisFlow := models.Flow{}

	redisFlowKey := builder.BuildRedisFlowKey(id)
	if err := r.client.Get(ctx, redisFlowKey).Scan(&redisFlow); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.Flow{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return dto.Flow{}, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Printf("redis flow: %v", redisFlow)

	flow, err := converter.ToFlowDTO(redisFlow)
	if err != nil {
		return dto.Flow{}, fmt.Errorf("%s: %w", op, err)
	}

	return flow, nil
}

func (r *Redis) SaveFlow(ctx context.Context, flow dto.Flow, ttl time.Duration) error {
	const op = "infrastructure.redis.SaveFlow"

	redisFlow := converter.ToFlowRedis(flow)

	redisFlowKey := builder.BuildRedisFlowKey(redisFlow.ID)
	if err := r.client.Set(ctx, redisFlowKey, redisFlow, ttl).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
