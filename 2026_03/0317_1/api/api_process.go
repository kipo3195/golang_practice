package api

import (
	"context"
	"math/rand"
	"sync"
	"test/entity"
	"time"
)

func APIProcess(ctx context.Context, wg *sync.WaitGroup, apiProcessNum int) chan entity.HealthStatus {

	apiStatusChan := make(chan entity.HealthStatus)

	go func(apiProcessNum int) {
		defer wg.Done()

		for i := 1; i <= apiProcessNum; i++ {
			r := rand.Intn(900) + 100
			delay := time.Duration(r) * time.Millisecond

			time.Sleep(delay)
			apiStatusChan <- entity.HealthStatus{
				ServiceName: "api",
				IsAlive:     true,
				Latency:     delay,
			}
		}

		close(apiStatusChan)

	}(apiProcessNum)

	return apiStatusChan
}
