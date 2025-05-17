package repository

import (
	"context"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/redis/go-redis/v9"
	"time"
)

type Storage interface {
	// Push adds a value to the storage with a specified key and expiration time.
	// Expiration is in seconds.
	Push(key string, value string, expiration time.Duration) error

	// Pop removes a value from the storage by its key and returns it.
	Pop(key string) (string, error)

	// Get retrieves a value from the storage by its key.
	// Without deleting it.
	Get(key string) (string, error)

	// Del removes a value from the storage by its key.
	Del(key string) error

	// Exists checks if a value exists in the storage by its key.
	Exists(key string) (bool, error)

	// Expire sets the expiration time for a value in the storage by its key.
	// Sets for the existing key.
	Expire(key string, expiration time.Duration) error

	// Clear removes all values from the storage.
	Clear() error
}

type RedisStorage struct {
	rdb *redis.Client
	ctx context.Context
}

func (r RedisStorage) Push(key string, value string, expiration time.Duration) error {
	return r.rdb.Set(r.ctx, key, value, expiration).Err()
}

func (r RedisStorage) Pop(key string) (string, error) {
	res, err := r.rdb.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}
	err = r.rdb.Del(r.ctx, key).Err()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r RedisStorage) Get(key string) (string, error) {
	logging.Logger.Debug("Getting value from redis with key: ", key)
	res, err := r.rdb.Get(r.ctx, key).Result()
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to get value from redis with key: ", key)
		return "", err
	}
	return res, nil
}

func (r RedisStorage) Del(key string) error {
	err := r.rdb.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisStorage) Exists(key string) (bool, error) {
	res, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r RedisStorage) Expire(key string, expiration time.Duration) error {
	err := r.rdb.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisStorage) Clear() error {
	err := r.rdb.FlushAll(r.ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisStorage(rdb *redis.Client, ctx context.Context) Storage {
	return &RedisStorage{rdb: rdb, ctx: ctx}
}
