package main

import (
	"context"
	"fmt"
	"sync"
	"test/producer"
	"test/workerpool"
)

// 특정 사용자에게 미확인 이벤트를 전송하는 워커풀 생성

// Sticky Task Scheduler문제
// 실시간 로그 프로세서를 구현
// 시스템의 효율성을 위해 특정 장치(Device)에서 오는 로그는 항상 동일한 워커가 처리하여 순차성과 로컬 캐시 효율을 높여야 합니다.
// 요구 사항Fixed Pool Size: 총 N개의 워커 고루틴이 상주하며 대기해야 합니다
// Stickiness: 각 작업은 DeviceID를 가집니다.
// 동일한 DeviceID를 가진 작업은 항상 동일한 워커 고루틴에 할당되어야 합니다
// Concurrency: 서로 다른 DeviceID를 가진 작업들은 병렬로 처리되어야 합니다
// Graceful Shutdown: 모든 작업이 완료되면 워커들을 안전하게 종료해야 합니다.

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// 이벤트 생성 로직 (producer)
	taskChan := producer.Produce()

	// N개의 워커
	workerNum := 5
	wp := workerpool.NewWorkerPool(workerNum)

	var wg sync.WaitGroup

	wp.Init(ctx, &wg)

	// task를 워커가 처리할 수 있도록
	for task := range taskChan {
		// 어떤 워커에 분배할지 hash값 구하기
		workerIdx := getWorkerHash(task.DeviceID, workerNum)
		// 분배
		wp.Workers[workerIdx].EventChan <- task
	}

	// 명시적 종료를 통한 graceful shutdown
	cancel()

	wg.Wait()
	fmt.Println("작업 종료.")

}

func getWorkerHash(deviceId int, workerNum int) int {
	return deviceId % workerNum
}
