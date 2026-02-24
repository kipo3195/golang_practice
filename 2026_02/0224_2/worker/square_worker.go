package worker

import "sync"

// 함수가 다른 함수를 반환하고 있습니다. 이를 고차 함수라고 부릅니다.
func SquareWorker() Work {
	// 이 함수는 클로저의 구조를 갖는다 (클로저의 특징인 외부 변수 캡쳐 로직 X)
	// 자신이 생성된 영역을 활용할 수 있는 함수 이기 때문이다.
	return func(n int) int {
		return n * n
	}
}

func SquareProcess(s []int, worker Work) []int {
	resultChan := make(chan int, len(s))
	var wg sync.WaitGroup

	// 슬라이스를 고루틴에 할당, worker에서 연산
	for _, v := range s {
		wg.Add(1)

		go func(target int) {

			defer wg.Done()
			resultChan <- worker(target)

		}(v)
	}

	go func() {
		// 모든 고루틴이 종료되고 난 다음 채널 종료처리
		wg.Wait()
		close(resultChan)
	}()

	var results []int
	for res := range resultChan {
		results = append(results, res)
	}

	return results
}
