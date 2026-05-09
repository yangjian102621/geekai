package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import "sync"

// UserLockManager 提供基于用户ID的TryLock功能，确保同一用户并发请求串行化
type UserLockManager struct {
	mu    sync.Mutex
	locks map[uint]bool
}

func NewUserLockManager() *UserLockManager {
	return &UserLockManager{mu: sync.Mutex{}, locks: make(map[uint]bool)}
}

// TryLock 尝试为指定用户加锁。若已被占用返回 false
func (m *UserLockManager) TryLock(userId uint) bool {
	if userId == 0 {
		return true
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.locks[userId] {
		return false
	}
	m.locks[userId] = true
	return true
}

// Unlock 释放指定用户的锁
func (m *UserLockManager) Unlock(userId uint) {
	if userId == 0 {
		return
	}
	m.mu.Lock()
	delete(m.locks, userId)
	m.mu.Unlock()
}
