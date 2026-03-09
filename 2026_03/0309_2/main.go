package main

import (
	"fmt"
	"sync"
	"time"
)

// https://gemini.google.com/share/4dcf34e47acd
// 문제 2: 워커 풀 패턴 구현 (Worker Pool with Channels)

func main() {

	jobs := make(chan int, 10)

	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	// 채널 닫기를 명시함으로써, for-range가 수행가능하도록
	close(jobs)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for value := range jobs {
				fmt.Printf("Worker %d started job %d\n", idx, value)
				time.Sleep(500 * time.Millisecond)
				fmt.Printf("Worker %d finished job %d\n", idx, value)
			}
		}(i)
	}

	// "작업 도중 에러가 발생하면 어떻게 전파하면 좋을까요?"
	// 답변 방향: 에러 전용 채널(errCh)을 따로 만들거나, golang.org/x/sync/errgroup 패키지를 사용하여 하나라도 에러가 나면 전체 작업을 중단하는 방식을 고려할 수 있습니다.

	wg.Wait()
	fmt.Println("작업 완료.")

}
