package main

import (
	"context"
	"fmt"
	"sync"
	"test/access"
	"test/chat"
	"time"
)

// 멀티 토픽 데이터 수집기 (Multi-Topic Aggregator)

// 당신은 메신저의 '채팅 로그'와 '접속 로그' 두 가지 데이터를 수집하여 통계를 내는 시스템을 만들고 있습니다.
// 각 로그는 별도의 고루틴에서 처리되며,
// 두 종류의 로그 수집이 모두 끝나는 순간 최종 리포트를 출력하고 프로그램을 종료해야 합니다.

//https://gemini.google.com/share/10634d3cd20d

// 데이터 생성: Chat 그룹은 3개의 메시지를, Access 그룹은 5개의 메시지를 생성합니다. (각기 다른 속도로 생성)
// 비동기 추적: sync.WaitGroup을 사용하여 두 그룹의 완료 상태를 추적하세요.
// 비동기 Wait 패턴: Wait()을 별도의 고루틴에서 실행하여 모든 작업이 끝나는 즉시 결과 채널(done)을 닫거나 신호를 보내세요.
// 최종 출력: main 함수는 다른 일을 하다가(예: 1초마다 "." 출력), done 신호를 받으면 "모든 로그 수집 완료, 리포트를 생성합니다"를 출력하고 종료해야 합니다.

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	report := make([]string, 0)

	wg.Add(2)
	chatChan := chat.ChatProduceProcess(ctx, &wg, 3)
	accessChan := access.AccessProduceProcess(ctx, &wg, 5)

	combinedChan := merge(chatChan, accessChan)

	resultChan := make(chan struct{})
	go func() {
		wg.Wait()
		// 리포트 생성
		close(resultChan)

	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println(".")
		case log := <-combinedChan:
			report = append(report, log)
		case <-resultChan:
			// 모든 생산 고루틴이 종료되었음이 보장된 시점
			// (채널이 닫혔으므로 남은 데이터가 있다면 마저 수집)
			fmt.Println("모든 로그 수집 완료, 리포트를 생성합니다.")

			// ⭐ 중요: merge 함수가 이미 모든 채널을 combinedChan으로 합쳤으므로
			// combinedChan에 남은 데이터만 마저 털어내면 됩니다.
			for len(combinedChan) > 0 {
				report = append(report, <-combinedChan)
			}

			fmt.Println("Final Report:", report)
			return
		}
	}
}

// 1. 채널 자체가 Thread-Safe (고루틴 안전) 합니다
// Go의 채널은 내부적으로 hchan이라는 구조체로 되어 있고, 데이터를 넣거나(send) 뺄 때(receive) 내부적으로 뮤텍스(Mutex) 잠금을 수행합니다.
// 여러 고루틴이 하나의 out 채널에 동시에 out <- n을 실행하더라도, Go 런타임이 차례대로 줄을 세워 데이터를 안전하게 전달합니다.
// 사용자는 별도의 sync.Mutex를 직접 구현할 필요가 없습니다.

// 2. 메모리 공유가 아닌 '데이터 전달'
// Race Condition은 보통 여러 스레드가 같은 메모리 공간(예: 전역 변수, 슬라이스)을 동시에 수정하려고 할 때 발생합니다.
// merge 구조에서는 고루틴들이 report라는 슬라이스를 직접 수정하는 것이 아니라, out이라는 통로에 데이터를 던지기만 합니다.
// 실제 슬라이스에 append를 하는 행위는 최종적으로 for msg := range combinedChan 루프가 도는 단 하나(Main)의 고루틴에서만 일어납니다. 즉, 쓰기 작업이 직렬화(Serialized)되므로 안전합니다.

// 3. 구조적 안전성 (Fan-In)
// 사용자가 우려하시는 Race Condition이 발생하려면 다음과 같은 상황이어야 합니다.
// 위험한 상황: 여러 고루틴이 report = append(report, msg)를 동시에 실행할 때 (이건 위험합니다!)
// 안전한 상황 (Merge 패턴): * 고루틴 A, B, C -> out 채널에 데이터 송신 (채널이 동기화 처리함)
// 메인 고루틴 -> out 채널에서 하나씩 꺼내서 report에 추가 (단일 쓰레드 작업)

func merge(cs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	output := func(c <-chan string) {
		defer wg.Done()
		for n := range c { // 채널이 닫힐 때까지 루프 (Zero Value 걱정 없음!)
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// 모든 채널이 소모되면 out 채널을 닫아줌
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
