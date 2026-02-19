package inMemory

import (
	"context"
	"log"
	"sync"
	"time"
)

// K는 비교 가능(comparable)해야 맵의 키가 될 수 있음
// V는 어떤 타입(any)이든 상관없음
type cacheInterfaceImpl[K comparable, V any] struct {
	data     map[K]item[V]
	mu       sync.RWMutex
	interval time.Duration
	cancel   context.CancelFunc
}

// 값과 만료 시간을 함께 관리할 구조체
type item[V any] struct {
	value V
	exp   time.Time
}

func NewCacheInterfaceImpl[K comparable, V any](interval time.Duration) cacheInterface[K, V] {

	// context 생성
	ctx, cancel := context.WithCancel(context.Background())

	c := &cacheInterfaceImpl[K, V]{
		data:     make(map[K]item[V]),
		interval: interval,
		cancel:   cancel,
	}

	// 비동기 고루틴
	go c.CleanUp(ctx)

	return c
}

func (r *cacheInterfaceImpl[K, V]) Set(key K, value V, exp time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = item[V]{
		value: value,
		exp:   time.Now().Add(exp), // 지금으로 부터 exp 이후
	}

}
func (r *cacheInterfaceImpl[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	// 지금이 만료보다 지났나?
	if time.Now().After(item.exp) {
		var zero V
		return zero, false
	}

	return item.value, true
}

func (r *cacheInterfaceImpl[K, V]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.data[key]
	if exists {
		delete(r.data, key)
	}
}

func (r *cacheInterfaceImpl[K, V]) CleanUp(ctx context.Context) {

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			for key, value := range r.data {
				exp := value.exp

				// 지금이 만료보다 지났나?
				if time.Now().After(exp) {
					// 지났으면 삭제
					delete(r.data, key)
				}
			}
			r.mu.Unlock()
		case <-ctx.Done():
			log.Println("clean up 종료.")
			return
		}
	}
}

func (r *cacheInterfaceImpl[K, V]) Close() {
	// graceful shutdown.
	r.cancel()
}
