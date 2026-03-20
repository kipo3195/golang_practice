package main

import (
	"context"
	"fmt"
	"sync"
	"test/produce"
	"test/workerpool"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageChan := produce.Producer(ctx, 100)

	workerpool := workerpool.NewWorkerPool(imageChan, 5)
	var wg sync.WaitGroup
	workerpool.Worker(ctx, &wg)

	// 모든 작업이 끝났다 (타임아웃, 작업완료)
	go func() {
		wg.Wait()
		// 채널 닫기
		close(workerpool.ResultChan)
	}()

	// 실시간 출력
	// ResultChan의 buffer가 없다면 워커들은 ResultChan에서 출력할때 까지 대기해야 함.
	// buffer를 주는 것을 통해 자기 일을 좀 더 빠르게 처리할 수 있음.
	// 실제로 70초반까지 출력되던 것을 워커의 수 만큼 buffer를 주니 80초반까지 작업을 수행함.
	for result := range workerpool.ResultChan {
		fmt.Printf("ID :%d, workingTime :%d, workingThread :%d \n", result.ID, result.WorkingTime, result.ThreadID)
	}

	fmt.Printf("작업 종료.")
}
