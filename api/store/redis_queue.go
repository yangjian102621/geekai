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
