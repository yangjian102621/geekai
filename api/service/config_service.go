package service

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Copyright 2023 The Geek-AI Authors. All rights reserved.
// Use of this source code is governed by a Apache-2.0 license
// that can be found in the LICENSE file.
// @Author yangjian102621@163.com
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/store/model"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ConfigService 统一的配置访问、缓存与通知服务
type ConfigService struct {
	db       *gorm.DB
	rdb      *redis.Client
	mu       sync.RWMutex
	cache    map[string]json.RawMessage
	watchers map[string][]chan struct{}
}

func NewConfigService(db *gorm.DB, rdb *redis.Client) *ConfigService {
	s := &ConfigService{
		db:       db,
		rdb:      rdb,
		cache:    make(map[string]json.RawMessage),
		watchers: make(map[string][]chan struct{}),
	}
	go s.subscribe()
	return s
}

// Get 以原始 JSON 获取配置（带本地缓存）
func (s *ConfigService) Get(key string) (json.RawMessage, error) {
	s.mu.RLock()
	if v, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return v, nil
	}
	s.mu.RUnlock()

	var cfg model.Config
	if err := s.db.Where("name", key).First(&cfg).Error; err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[key] = json.RawMessage(cfg.Value)
	s.mu.Unlock()
	return json.RawMessage(cfg.Value), nil
}

// GetInto 将配置解析进传入结构体
func (s *ConfigService) GetInto(key string, dest interface{}) error {
	data, err := s.Get(key)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, dest)
}

// Set 设置配置并写入数据库，同时触发通知
func (s *ConfigService) Set(key string, config json.RawMessage) error {
	value := string(config)
	cfg := model.Config{Name: key, Value: value}
	if err := s.db.Where("name", key).FirstOrCreate(&cfg, model.Config{Name: key}).Error; err != nil {
		return err
	}
	if cfg.Id > 0 {
		cfg.Value = value
		if err := s.db.Updates(&cfg).Error; err != nil {
			return err
		}
	}
	s.mu.Lock()
	s.cache[key] = json.RawMessage(value)
	s.mu.Unlock()
	s.notifyLocal(key)
	s.publish(key)
	return nil
}

// Watch 返回一个通道，当指定 key 发生变化时收到事件
func (s *ConfigService) Watch(key string) <-chan struct{} {
	ch := make(chan struct{}, 1)
	s.mu.Lock()
	s.watchers[key] = append(s.watchers[key], ch)
	s.mu.Unlock()
	return ch
}

func (s *ConfigService) notifyLocal(key string) {
	s.mu.RLock()
	list := s.watchers[key]
	s.mu.RUnlock()
	for _, ch := range list {
		select { // 非阻塞通知
		case ch <- struct{}{}:
		default:
		}
	}
}

// 通过 Redis 发布配置变更，便于多实例同步
func (s *ConfigService) publish(key string) {
	if s.rdb == nil {
		return
	}
	channel := "config:changed"
	if err := s.rdb.Publish(context.Background(), channel, key).Err(); err != nil {
		logger.Warnf("publish config change failed: %v", err)
	}
}

func (s *ConfigService) subscribe() {
	if s.rdb == nil {
		return
	}
	channel := "config:changed"
	sub := s.rdb.Subscribe(context.Background(), channel)
	for msg := range sub.Channel() {
		key := msg.Payload
		logger.Infof("config changed: %s", key)
		// 失效本地缓存并本地广播
		s.mu.Lock()
		delete(s.cache, key)
		s.mu.Unlock()
		s.notifyLocal(key)
	}
}

// Test 预留统一测试入口，根据 key 执行连通性检查
func (s *ConfigService) Test(key string) (string, error) {
	// TODO: 实现各配置类型的测试逻辑
	return fmt.Sprintf("%s ok", key), nil
}
