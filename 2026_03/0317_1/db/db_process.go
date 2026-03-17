package db

import (
	"context"
	"sync"
	"test/entity"
	"time"
)

func DBProcess(ctx context.Context, wg *sync.WaitGroup, dbProcessNum int) chan entity.HealthStatus {

	dbStatusChan := make(chan entity.HealthStatus)

	go func(apiProcessNum int) {
		defer wg.Done()
		ticker := time.NewTicker(400 * time.Millisecond)
		defer ticker.Stop()
		for i := 1; i <= apiProcessNum; i++ {
			<-ticker.C
			dbStatusChan <- entity.HealthStatus{
				ServiceName: "db",
				IsAlive:     true,
				Latency:     400,
			}
		}

		close(dbStatusChan)

	}(dbProcessNum)

	return dbStatusChan
}
