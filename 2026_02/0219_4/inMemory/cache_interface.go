package inMemory

import "context"

type cacheInterface interface {
	Set(key string, value string, ttl string)
	Get(key string) string
	Delete(key string)
	CleanUp(ctx context.Context, interval int)
}
