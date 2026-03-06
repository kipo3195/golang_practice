package main

import (
	"context"
	"fmt"
	"test/collector"
	"test/dispatcher"
	"test/entity"
)

// 외부 알림 API는 초당 처리량(Rate Limit)이 제한되어 있어,
// 무작정 고루틴을 생성하면 오류가 발생합니다.
// 요구사항Worker Pool 구현: 최대 $N$개의 고루틴만 동시에 작동하도록 제한하는 Dispatcher를 구현하세요.
// 타임아웃 처리: 각 메시지 발송 작업은 최대 500ms 안에 완료되어야 하며, 초과 시 해당 작업은 실패로 간주하고 다음으로 넘어갑니다.
// 결과 수집: 모든 작업이 끝난 후, 성공한 개수와 실패한 개수를 안전하게(Thread-safe) 집계하여 출력하세요.
// Graceful Shutdown: 모든 작업이 전달된 후, 현재 실행 중인 고루틴들이 모두 종료될 때까지 기다렸다가 프로그램을 종료하세요.

// https://gemini.google.com/share/53ae569a4620

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 버퍼가 없다면 채널에 데이터를 하나씩 넣고 뺄때까지 기다렸다가 하나 넣고 하는 구조.
	// 효율적인 측면에서 workerpool를 사용하는건데 하나씩 빼갈때까지 collector를 쉬게 할 순 없다.
	// worker들이 기아상태에 빠지지 않게 worker의 수 *2 만큼을 버퍼의 사이즈로 잡는다.
	// 다만, 버퍼 사이즈를 너무 크게 잡는다면 collector와 worker사이의 병목이 줄어들어 속도는 빨라지지만, 그만큼 메모리를 점유하게 된다.
	// 그러므로 시스템 가용 메모리와 메시지의 크기를 고려해서 적절한 타협점을 찾는 것이 중요하다.
	msgChan := make(chan entity.Message, 6)
	collector := collector.NewCollector(msgChan)

	// 메시지
	msg := []string{"apple", "banana"}
	collector.Collect(ctx, msg)

	resultChan := make(chan entity.Result)
	dispatcher := dispatcher.NewDispatcher(3, msgChan, resultChan)

	dispatcher.Work(ctx)

	// resultChan에서 데이터를 뽑는 처리를 하는 영역
	var success int
	var fail int
	for value := range resultChan {
		if value.Success {
			success++
		} else {
			fail++
		}
	}

	fmt.Printf("결과 성공 : %d, 실패 : %d", success, fail)
}
