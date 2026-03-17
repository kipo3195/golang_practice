package redis

import (
	"context"
	"sync"
	"test/entity"
	"time"
)

func RedisProcess(ctx context.Context, wg *sync.WaitGroup, redisProcessNum int) chan entity.HealthStatus {

	redisStatusChan := make(chan entity.HealthStatus)

	go func(redisProcessNum int) {
		defer wg.Done()
		ticker := time.NewTicker(600 * time.Millisecond)
		defer ticker.Stop()
		for i := 1; i <= redisProcessNum; i++ {
			<-ticker.C
			redisStatusChan <- entity.HealthStatus{
				ServiceName: "redis",
				IsAlive:     true,
				Latency:     600,
			}
		}

		close(redisStatusChan)

	}(redisProcessNum)

	return redisStatusChan
}
