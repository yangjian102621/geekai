package xxl

import "sync"

//任务列表 [JobID]执行函数,并行执行时[+LogID]
type taskList struct {
	mu   sync.RWMutex
	data map[string]*Task
}

// Set 设置数据
func (t *taskList) Set(key string, val *Task) {
	t.mu.Lock()
	t.data[key] = val
	t.mu.Unlock()
}

// Get 获取数据
func (t *taskList) Get(key string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[key]
}

// GetAll 获取数据
func (t *taskList) GetAll() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

// Del 设置数据
func (t *taskList) Del(key string) {
	t.mu.Lock()
	delete(t.data, key)
	t.mu.Unlock()
}

// Len 长度
func (t *taskList) Len() int {
	return len(t.data)
}

// Exists Key是否存在
func (t *taskList) Exists(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.data[key]
	return ok
}
