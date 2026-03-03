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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 채널은 넣는 로직이 있다면 빼는로직도 반드시 있어야함.
	resultChan := make(chan dto.Result, len(jobs))

	var wg sync.WaitGroup
	for i := 0; i < len(jobs); i++ {
		wg.Add(1)
		go worker.Worker(ctx, &wg, jobs[i], resultChan)
	}

	// 모든 고루틴이 종료 = 채널 닫기
	// 기다리고 닫는 역할을 별도 고루틴에 위임 (Non-blocking)
	go func() {
		wg.Wait()
		// 모든 고루틴이 종료되면 wg.Wait이 수행될꺼고.. 그러면 채널이 닫힘 그러니까 아래 for문에서 채널이 닫힐때까지를 명시하는 for range로 처리
		close(resultChan)

	}()

	for res := range resultChan {
		log.Printf("결과 수신: %+v", res)
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("⚠️ 전체 작업 시간 초과로 일부 작업이 취소되었습니다.")
	}

	log.Println("작업종료")

}
