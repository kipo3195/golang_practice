package access

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func AccessProduceProcess(ctx context.Context, wg *sync.WaitGroup, processNum int) chan string {

	accessChan := make(chan string)
	ticker := time.NewTicker(4 * time.Second)

	go func(processNum int) {
		defer wg.Done()
		defer close(accessChan)
		defer ticker.Stop()

		for i := 1; i <= processNum; i++ {
			// 설정된 시간 만큼 대기
			<-ticker.C
			fmt.Println("[Access] logging")
			accessChan <- fmt.Sprintf("access log %d", i)
		}

	}(processNum)

	return accessChan
}
