package repository

import (
	"context"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisRepository struct {
	client *redis.Client
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	logging.Logger.Infof("Getting key: %v", key)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		logging.Logger.WithError(err).Errorf("Failed to get key: %s", key)
		return "", err
	}
	return val, nil
}

func (r *redisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	logging.Logger.Infof("Setting key: %s with value: %s and ttl: %v", key, value, ttl)
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		logging.Logger.WithError(err).Errorf("Failed to set key: %s with value: %s", key, value)
		return err
	}
	return nil
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{
		client: client,
	}
}
