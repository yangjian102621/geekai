package types

import (
	"sync"
)

type LockedMap struct {
	lock sync.RWMutex
	data map[string]interface{}
}

func NewLockedMap() *LockedMap {
	return &LockedMap{
		lock: sync.RWMutex{},
		data: make(map[string]interface{}),
	}
}

func (m *LockedMap) Put(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = value
}

func (m *LockedMap) Get(key string) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.data[key]
}

func (m *LockedMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.data, key)
}

func (m *LockedMap) ToList() []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()

	var s = make([]interface{}, 0)
	for _, v := range m.data {
		s = append(s, v)
	}
	return s
}
