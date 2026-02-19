package inMemory

import (
	"context"
	"time"
)

type cacheInterface[K comparable, V any] interface {
	Set(key K, value V, exp time.Duration)
	Get(key K) (V, bool)
	Delete(key K)
	CleanUp(ctx context.Context)
	Close()
}
