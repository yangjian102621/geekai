package service

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 雪花算法实现
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int
	sequence      int
}

func NewSnowflake() *Snowflake {
	return &Snowflake{
		lastTimestamp: -1,
		workerID:      0, // TODO: 增加 WorkID 参数
		sequence:      0,
	}
}

// Next 生成一个新的唯一ID
func (s *Snowflake) Next(raw bool) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000 // 转换为毫秒
	if timestamp < s.lastTimestamp {
		return "", fmt.Errorf("clock moved backwards. Refusing to generate id for %d milliseconds", s.lastTimestamp-timestamp)
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & 4095
		if s.sequence == 0 {
			timestamp = s.waitNextMillis()
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp
	id := (timestamp << 22) | (int64(s.workerID) << 10) | int64(s.sequence)
	if raw {
		return fmt.Sprintf("%d", id), nil
	}
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%d", now.Year(), now.Month(), now.Day(), id), nil
}

func (s *Snowflake) waitNextMillis() int64 {
	timestamp := time.Now().UnixNano() / 1000000
	for timestamp <= s.lastTimestamp {
		timestamp = time.Now().UnixNano() / 1000000
	}
	return timestamp
}
