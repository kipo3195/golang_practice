package main

import (
	"context"
	"log"
	"sync"
	"test/dto"
	"test/worker"
	"time"
)

// https://gemini.google.com/share/6b6f8c93d56

func main() {
	jobs := []string{"job1", "job2", "fail", "job3", "slow_job"}

	// TODO: 1. 500ms 타임아웃 컨텍스트 생성
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)

	resultChan := make(chan dto.Result)
	defer close(resultChan)

	var wg sync.WaitGroup
	for i := 0; i < len(jobs); i++ {
		wg.Add(1)
		go worker.Worker(ctx, &wg, jobs[i], resultChan)
	}

	wg.Wait()

	log.Println("작업종료")

}
