package store

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"geekai/utils"

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

func (q *RedisQueue) RPush(value any) error {
	return q.client.RPush(q.ctx, q.name, utils.JsonEncode(value)).Err()
}

func (q *RedisQueue) LPush(value any) error {
	return q.client.LPush(q.ctx, q.name, utils.JsonEncode(value)).Err()
}

func (q *RedisQueue) LPop(value any) error {
	result, err := q.client.BLPop(q.ctx, 0, q.name).Result()
	if err != nil {
		return err
	}
	return utils.JsonDecode(result[1], value)
}

func (q *RedisQueue) RPop(value any) error {
	result, err := q.client.BRPop(q.ctx, 0, q.name).Result()
	if err != nil {
		return err
	}
	return utils.JsonDecode(result[1], value)
}

func (q *RedisQueue) Size() (int64, error) {
	return q.client.LLen(q.ctx, q.name).Result()
}

func (q *RedisQueue) Clear() error {
	return q.client.Del(q.ctx, q.name).Err()
}
