package store

import (
	"chatplus/core/types"
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(config *types.AppConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Url(),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
