package main

import (
	"context"
	"log"
	"sync"
	"test/worker"
	"time"
)

// 당신은 결제 시스템의 검증 로직을 담당합니다.
// 결제가 승인되려면 4가지 검증(신원 확인, 잔액 조회, 사기 탐지, 블랙리스트 체크)을 동시에 통과해야 합니다.
// 하나라도 실패 하면 나머지 작업도 모두 취소

// 각 검증을 하나의 고루틴이 처리한다
// 내부 로직에서 실패 -> cancle()을 호출하고 고루틴 종료하기
// 실패 트리거는 timeout?
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	auth := []string{"idCheck", "exchangeCheck", "lierCheck", "blacklistCheck"}

	var wg sync.WaitGroup

	resultChan := make(chan string, len(auth))

	for _, a := range auth {
		log.Println(a)
		wg.Add(1)
		// 얘는 클로저가 아니지 않음? a를 계속 재할당하니까
		go worker.Worker(ctx, cancel, &wg, resultChan, a)
	}

	// 고루틴의 종료도 고루틴에게 맞김 (블로킹 방지 )
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 길이는 0, 용량은 auth와 같음
	result := make([]string, 0, len(auth))
	// 기다려주는놈이 필요함.
	for r := range resultChan {
		result = append(result, r)
	}

	if ctx.Err() == context.Canceled {
		log.Println("누군가 취소를 눌렸음")
	} else if ctx.Err() == context.DeadlineExceeded {
		log.Println("시간초과!")
	} else {
		log.Println("작업 종료. result :", result)
	}

}
