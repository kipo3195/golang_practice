package storage

import "sync"

type Storage struct {
	TotalCount int64
	mu         sync.RWMutex // 포인터로 정의하게 되면 init 필요. mutex는 일반적으로 값타입으로 선헌하는 것이 일반적
}

func NewStorage() *Storage {
	return &Storage{}
}

func (r *Storage) Add(count int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.TotalCount += count
}

func (r *Storage) Result() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.TotalCount
}
