package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"sync"
)

type MKey interface {
	string | int | uint
}
type MValue interface {
	*WsClient | *ChatSession | context.CancelFunc | []Message
}
type LMap[K MKey, T MValue] struct {
	lock sync.RWMutex
	data map[K]T
}

func NewLMap[K MKey, T MValue]() *LMap[K, T] {
	return &LMap[K, T]{
		lock: sync.RWMutex{},
		data: make(map[K]T),
	}
}

func (m *LMap[K, T]) Put(key K, value T) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = value
}

func (m *LMap[K, T]) Get(key K) T {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.data[key]
}

func (m *LMap[K, T]) Has(key K) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.data[key]
	return ok
}

func (m *LMap[K, T]) Delete(key K) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.data, key)
}

func (m *LMap[K, T]) ToList() []T {
	m.lock.Lock()
	defer m.lock.Unlock()

	var s = make([]T, 0)
	for _, v := range m.data {
		s = append(s, v)
	}
	return s
}
