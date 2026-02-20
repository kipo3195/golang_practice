package main

import (
	"0220_2/storage"
	"log"
	"sync"
)

// 글자수의 합을 구하기
// 5개의 고루틴 사용

func main() {

	logs := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
	var totalLength int64 // 결과값이 저장될 공유 변수

	jobs := make(chan string, len(logs))
	var wg sync.WaitGroup

	workerPoolNum := 5

	// 결과 저장을 위한 공통 storage
	storage := storage.NewStorage()

	// worker init
	for i := 1; i <= workerPoolNum; i++ {
		wg.Add(1)
		go counter(&wg, jobs, storage)
	}

	// 채널에 데이터 주입
	for _, value := range logs {
		jobs <- value
	}
	// 채널 닫기 = conter 내부에서 종료처리 할 수 있도록 ok를 false로 바꾸는 역할
	close(jobs)

	// 끝날때 까지 대기
	wg.Wait()

	// 결과 가져오기
	totalLength = storage.Result()

	log.Printf("작업 완료. 결과 :%d", totalLength)

}

func counter(wg *sync.WaitGroup, jobs <-chan string, storage *storage.Storage) {
	defer wg.Done()
	var innerCount int64

	// for - select 방식
	// for {
	// 	select {
	// 	case log, ok := <-jobs:
	// 		if !ok {
	// 			// 채널 닫힘
	// 			storage.Add(innerCount)
	// 			return
	// 		} else {
	// 			innerCount = (int64)(len(log))
	// 		}
	// 	}
	// }

	// for - range 방식
	for log := range jobs {
		innerCount += (int64)(len(log))
	}

	// 채널 닫힌 후 결과 저장
	storage.Add(innerCount)
}
