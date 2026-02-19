package inMemory

import (
	"context"
	"sync"
	"time"
)

type cacheInterfaceImpl struct {
	cacheMap map[string]string
	ttlMap   map[string]time.Time
	mu       sync.RWMutex
}

func NewCacheInterfaceImpl() cacheInterface {
	return &cacheInterfaceImpl{
		cacheMap: make(map[string]string),
		ttlMap:   make(map[string]time.Time),
	}
}

func (r *cacheInterfaceImpl) Set(key string, value string, ttl string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cacheMap[key] = value
	layout := "20060102150405"
	// 넣을때 부터 파싱해서 시간으로 넣고 꺼내면 비교만..
	t, _ := time.Parse(layout, ttl)
	r.ttlMap[key] = t
}

func (r *cacheInterfaceImpl) Get(key string) string {
	r.mu.RLock()
	defer r.mu.Unlock()

	ttl, exists := r.ttlMap[key]
	if !exists {
		return ""
	}

	nowDate := time.Now()

	if ttl.After(nowDate) {
		// ttl이 더 나중 = 지나지 않았음
		return r.cacheMap[key]
	} else {
		// r.Delete(key) 메소드 호출은 데드락 유발..

		// RLock은 쓰기 lock인데, map의 데이터를 delete하게 되면 Runtime Panic 발생
		// delete(r.cacheMap, key)
		// delete(r.ttlMap, key)
	}

	return ""
}

func (r *cacheInterfaceImpl) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.cacheMap[key]
	if exists {
		delete(r.cacheMap, key)
		delete(r.ttlMap, key)
	}
}

func (r *cacheInterfaceImpl) CleanUp(ctx context.Context, interval int) {

	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nowDate := time.Now()
			r.mu.Lock()
			for key, ttl := range r.ttlMap {
				if ttl.Before(nowDate) {
					delete(r.cacheMap, key)
					delete(r.ttlMap, key)
				}
			}
			r.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}

}
