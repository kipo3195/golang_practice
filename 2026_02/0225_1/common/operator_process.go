package common

import "sync"

// 그 결과물들을 채널을 통해 수집하여 다시 슬라이스로 반환하는 함수를 작성하세요.
func OperatorProcess(s []int, w Operator) []int {

	// 고루틴 및 채널 활용: 생성된 작업을 여러 개의 고루틴으로 동시에 실행하고,
	// 결과 슬라이스의 순서는 상관없습니다.

	// 모든 고루틴이 종료됨을 감지하기 위한 장치
	var wg sync.WaitGroup

	// 결과를 담을 채널
	resultChan := make(chan int, len(s))

	for _, value := range s {

		// 고루틴 실행
		wg.Add(1)

		go func(target int) {
			// 고루틴 함수 종료시
			defer wg.Done()

			// 자체가 함수이므로..
			resultChan <- w(target)

		}(value)
	}

	// 동기화: 모든 고루틴이 종료될 때까지 기다린 후 채널을 닫아야 하며,
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var result []int
	// for range chan은 무조건 닫히는 것이 보장된 채널에 사용해야함. 아니면 deadlock
	for value := range resultChan {
		result = append(result, value)
	}

	return result
}
