package main

import (
	"log"
	"sort"
	"sync"
	"test/common"
	"test/sum"
)

// 클로저 활용: 정수 슬라이스를 받아 각 요소에 +1 을 수행하는 함수를 반환하는 '작업 생성기' 함수를 작성하세요.
// 고루틴 및 채널 활용: 생성된 작업을 여러 개의 고루틴으로 동시에 실행하고, 그 결과물들을 채널을 통해 수집하여 다시 슬라이스로 반환하는 함수를 작성하세요.
// 동기화: 모든 고루틴이 종료될 때까지 기다린 후 채널을 닫아야 하며, 결과 슬라이스의 순서는 상관없습니다.
// 테스트 코드: 작성한 로직이 올바르게 동작하는지 검증하는 단위 테스트(go test) 코드를 포함하세요.

func main() {

	var wg sync.WaitGroup

	// 어떤 데이터
	var s []int
	for i := 1; i <= 10; i++ {
		s = append(s, i)
	}

	// 어떤 동작
	op := sum.SumWorkOperator()

	// worker들이 데이터를 공유할 채널
	dataChan := make(chan int, 1000)

	// worker들이 결과값을 전달할 채널
	resultChan := make(chan int, 1000)

	// workerpool
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		w := common.NewWorker(dataChan, op)
		go w.Process(&wg, resultChan)
	}

	// 데이터 전송 로직
	for i := 1; i < 1000; i++ {
		dataChan <- i
	}

	// 고루틴 데이터 수신 채널 닫기 - range문은 채널이 종료됨을 가정한 로직이므로 Process에서도 자연스럽게 for루프 빠져나옴
	close(dataChan)

	go func() {
		// Process 함수에서 wg.Done()이 defer로 호출되었다 = recvChan이 닫혔다 = for 루프 빠져나왔다.
		wg.Wait()
		close(resultChan)
	}()

	var result []int
	for v := range resultChan {
		result = append(result, v)
	}

	sort.Ints(result)
	log.Println("결과 : ", result)

}
