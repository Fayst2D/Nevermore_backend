package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Set(ctx context.Context, key string, value string, time time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type repo struct {
	client *redis.Client
}

func New(cfg Config) (Storage, error) {
	redisURL := fmt.Sprintf("redis://%s:%s@%s/%d",
		cfg.User,
		cfg.Password,
		cfg.Url,
		cfg.DB,
	)
	fmt.Println("Redis URL:", redisURL)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	repo := &repo{
		client: redis.NewClient(opt),
	}

	return repo, nil
}

func (r *repo) Set(ctx context.Context, key string, value string, time time.Duration) error {
	err := r.client.Set(ctx, key, value, time).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in redis: %w", err)
	}

	return nil
}

func (r *repo) Get(ctx context.Context, key string) (string, error) {

	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key does not exist: %w", err)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key from redis: %w", err)
	}

	return data, nil
}

func (r *repo) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("failed to delete key from redis: %w", err)
	}
	return nil
}
