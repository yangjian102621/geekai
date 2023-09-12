package store

import (
	"chatplus/utils"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	name   string
	client *redis.Client
	ctx    context.Context
}

func NewRedisQueue(name string, client *redis.Client) *RedisQueue {
	return &RedisQueue{name: name, client: client, ctx: context.Background()}
}

func (q *RedisQueue) RPush(value interface{}) {
	q.client.RPush(q.ctx, q.name, utils.JsonEncode(value))
}

func (q *RedisQueue) LPush(value interface{}) {
	q.client.LPush(q.ctx, q.name, utils.JsonEncode(value))
}

func (q *RedisQueue) LPop(value interface{}) error {
	result, err := q.client.BLPop(q.ctx, 0, q.name).Result()
	if err != nil {
		return err
	}
	return utils.JsonDecode(result[1], value)
}

func (q *RedisQueue) RPop(value interface{}) error {
	result, err := q.client.BRPop(q.ctx, 0, q.name).Result()
	if err != nil {
		return err
	}
	return utils.JsonDecode(result[1], value)
}
