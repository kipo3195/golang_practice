package main

import (
	"sync"
	"time"
)

// 멀티 토픽 데이터 수집기 (Multi-Topic Aggregator)

// 당신은 메신저의 '채팅 로그'와 '접속 로그' 두 가지 데이터를 수집하여 통계를 내는 시스템을 만들고 있습니다.
// 각 로그는 별도의 고루틴에서 처리되며, 두 종류의 로그 수집이 모두 끝나는 순간 최종 리포트를 출력하고 프로그램을 종료해야 합니다.

//https://gemini.google.com/share/10634d3cd20d

func main() {

	resultChan := make(chan int)

	chatChan := make(chan int, 3)
	var chatWg sync.WaitGroup
	chatTicker := time.NewTicker(3 * time.Second)
	defer chatTicker.Stop()

	go func() {
		defer chatWg.Done()
		for i := 1; i <= 3; i++ {
			select {
			case <-chatTicker.C:
				resultChan <- i
			}
		}
	}()

	var noteWg sync.WaitGroup
	noteChan := make(chan int, 5)
	noteTicker := time.NewTicker(3 * time.Second)

	go func() {
		defer noteWg.Done()
		for i := 1; i <= 3; i++ {
			select {
			case <-chatTicker.C:
				resultChan <- i
			}
		}
	}()
	defer noteTicker.Stop()

}
