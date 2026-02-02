package redis

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
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

func (r *Redis) SaveFlow(ctx context.Context, flow dto.Flow, ttl time.Duration) error {
	redisFlow := converter.ToFlowRedis(flow)

	redisFlowKey := redisFlow.RedisKey()

	redisFlowBinary, err := redisFlow.MarshalBinary()
	if err != nil {
		return fmt.Errorf("%w: %w", infrastructure.ErrMarshalData, err)
	}

	if err = r.client.Set(ctx, redisFlowKey, redisFlowBinary, ttl).Err(); err != nil {
		return fmt.Errorf("%s: %w", "set flow to redis", err)
	}

	return nil
}
