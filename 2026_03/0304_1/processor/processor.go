package processor

import (
	"context"
	"sync"
	"time"
)

func Processor(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// in의 모든 데이터를 소진하고나면 for루프 종료
			// wg.Done()
			for v := range in {
				time.Sleep(10 * time.Millisecond)
				sqaure := v * v
				out <- sqaure
			}
		}()
	}

	// 모든 wg.Done()이 호출되고나면 채널 닫기
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
