package process

import (
	"context"
	"math/rand"
	"sync"
	"test/entity"
	"time"
)

func Process(ctx context.Context, wg *sync.WaitGroup, service string, processCnt int) chan entity.Status {

	statusChan := make(chan entity.Status)

	go func(processCnt int) {
		defer close(statusChan)
		defer wg.Done()

		for i := 1; i <= processCnt; i++ {
			r := rand.Intn(900) + 100
			time.Sleep(time.Duration(r) * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			case statusChan <- entity.Status{
				Service:   service,
				Latency:   time.Duration(r),
				Timestamp: time.Now(),
			}:
			}
		}

	}(processCnt)

	return statusChan
}
