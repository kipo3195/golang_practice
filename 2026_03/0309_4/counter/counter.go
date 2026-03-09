package counter

import (
	"sync/atomic"
)

type Counter struct {
	value int64
	//mu    sync.RWMutex
}

func NewCounter(value int64) *Counter {
	return &Counter{
		value: value,
	}
}

func (r *Counter) AddCounter() {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	// r.value = r.value + 1

	// sync/atomic 패키지의 AddInt64 등을 사용하여 하드웨어 수준에서 안전하게 연산하기 = mutex 필요 없음
	atomic.AddInt64(&r.value, 1)

}

func (r *Counter) GetCounter() int64 {
	// r.mu.RLocker()
	// defer r.mu.RLocker().Unlock()

	atomic.LoadInt64(&r.value)
	return r.value
}
