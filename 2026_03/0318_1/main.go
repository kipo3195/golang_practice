package main

import (
	"context"
	"fmt"
	"sync"
	"test/entity"
	"test/process"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)
	authChan := process.Process(ctx, &wg, "auth", 10)
	inventoryChan := process.Process(ctx, &wg, "inventory", 10)
	shippingChan := process.Process(ctx, &wg, "shipping", 10)

	combinedChan := merge(ctx, authChan, inventoryChan, shippingChan)

	// 현재까지 모인 데이터
	report := make(map[string]entity.Report)

ExitLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("time out!")
			// 여기서 탈출시 아래 wg.Wait()은 생산자들이 타임아웃 신호를 받고 모든 작업을 정리(wg.Done)할 때까지 메인 고루틴이 잠시 기다려주는 역할
			break ExitLoop
		case status, ok := <-combinedChan:
			if !ok {
				// combinedChan이 닫힘
				// 여기서 탈출시 모든 채널이 닫혔다는 의미 = 작업이 끝났다
				break ExitLoop
			}
			service, exists := report[status.Service]
			if !exists {
				report[status.Service] = entity.Report{}
			}

			// 집계
			service.Total += int(status.Latency)
			if status.Latency > 800 {
				fmt.Printf("%s 서비스 지연 발생! %dms", status.Service, status.Latency)
				service.LateCount++
			}
			service.Count++

			report[status.Service] = service
			fmt.Println("status : ", status)

		}
	}

	// 채널들이 모두 데이터를 전송할 때까지 기다림 (채널이 닫혔지만 combined채널이 아직 안닫힌 타이밍이 존재하지 않을까?)
	// 만약 위 for select를 별도 고루틴으로 정의했다고하면 wg.Wait()이 끝났다고 해도, 집계 고루틴이 맵에 값을 쓰고 있는 찰나에 메인 고루틴이 맵을 읽으려고 하면
	// fatal error: concurrent map iteration and map write**가 발생하며 프로그램이 터질 수 있습니다.
	// 그러므로 동기식으로 처리할 수 있도록 수정.
	wg.Wait()

	// 결국 map는 thread safe 하지 않으므로 쓰기와 닫기를 하나의 고루틴이 처리할 수 있어야함.
	for key := range report {
		s := report[key]
		fmt.Printf("집계 service :%s, avg :%d, late :%d\n", key, s.Total/s.Count, s.LateCount)
	}

}

func merge(ctx context.Context, cs ...chan entity.Status) chan entity.Status {

	combinedChan := make(chan entity.Status)

	// 내부에서 사용할 wg
	var innerWg sync.WaitGroup

	innerWg.Add(len(cs))

	// combinedChan에 데이터 in
	work := func(innerWg *sync.WaitGroup, c chan entity.Status) {
		defer innerWg.Done()

		for status := range c {
			select {
			case <-ctx.Done():
				return
			case combinedChan <- status:
			}
		}
	}

	// 실행
	for _, c := range cs {
		go work(&innerWg, c)
	}

	// 종료
	// innerWg.Done() 수행 보장 = 타임아웃 또는 모든 데이터 처리 완료
	go func() {
		innerWg.Wait()
		close(combinedChan)
	}()

	return combinedChan
}
