package main

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const (
	keyTTL = time.Hour * 24
)

type Redis struct {
	redis *redis.Client
}

type Cache interface {
	InsertURL(key, url string) (string, error)
	GetURL(key string) (string, error)
	DeleteURL(key string) error

	Ping() error
}

func NewCache() (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &Redis{redis: client}, nil
}

func (r *Redis) InsertURL(key, url string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return r.redis.Set(ctx, key, url, keyTTL).Result()
}

func (r *Redis) GetURL(key string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return r.redis.Get(ctx, key).Result()
}

func (r *Redis) DeleteURL(key string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return r.redis.Del(ctx, key).Err()
}

func (r *Redis) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := r.redis.Ping(ctx).Result()

	if err != nil || result != "PONG" {
		return fmt.Errorf("failed to ping redis: %v", err)
	}

	return nil
}
