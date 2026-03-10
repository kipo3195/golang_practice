package main

import (
	"context"
	"fmt"
	"sync"
	"test/producer"
	"test/workerpool"
	"time"
)

// Producer (생성자): 1부터 100까지의 숫자를 로그 메시지(int)로 생성하여 채널에 보냅니다.

// Worker Pool (워커 풀): * 최대 5개의 워커가 동시에 실행되어야 합니다.
// 각 워커는 채널에서 받은 숫자에 2를 곱하는 시뮬레이션 작업을 수행합니다.
// 작업 효율을 위해 각 작업당 100ms의 지연 시간(time.Sleep)이 발생한다고 가정합니다.

// Result Aggregator (소비자): 워커들이 처리한 모든 결과값의 총합을 계산하여 출력합니다.

// Context (제어):
// 전체 프로세스는 3초 이내에 완료되어야 합니다.
// 만약 3초가 지나면 현재 진행 중인 모든 고루틴은 안전하게 종료(Graceful Shutdown)되어야 하며, 그때까지 처리된 합계만 출력합니다.

// Synchronization: 모든 고루틴이 종료될 때까지 main 함수가 기다려야 합니다.

func main() {

	// 프로세스 제어를 위한 context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 모든 고루틴이 종료될때까지 기다리기 위해
	var wg sync.WaitGroup

	// 연산 데이터 전송용 채널
	dataChan := make(chan int)
	producer.Produce(ctx, dataChan, 100)

	// 결과 전송용 채널
	resultChan := make(chan int)

	workerpool := workerpool.NewWorkerPool(5)
	workerpool.Work(ctx, &wg, dataChan, resultChan)

	go func() {
		// 모든 고루틴 종료
		wg.Wait()
		close(resultChan) // 모든 고루틴 종료 채널을 닫아야 panic 발생하지 않음.
	}()

	// 누적 합계를 위한 변수
	totalSum := 0

	// ctx에 의해 타임아웃 -> 고루틴 종료 -> wg.Wait() -> close(resultChan)이므로
	for value := range resultChan {
		totalSum += value
	}

	fmt.Println("totalSum:", totalSum)
	fmt.Println("작업 종료.")

}
