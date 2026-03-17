package main

import (
	"context"
	"fmt"
	"sync"
	"test/api"
	"test/db"
	"test/entity"
	"test/redis"
)

// 1. 멀티 서버 헬스체크 모니터링
// 여러 대의 서버(API 서버, DB, Redis 등)에 상태 확인 요청을 보낸 뒤, 가장 먼저 응답이 오는 순서대로 대시보드에 표시해야 하는 상황입니다.
// 필요성: 서버마다 네트워크 상태가 다르므로, 동기적으로 하나씩 기다리면 전체 응답 시간이 매우 길어집니다.
// 구조: 각 서버 체크 고루틴이 결과를 채널로 던지고, 이를 merge하여 메인 루프에서 실시간으로 출력합니다.

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)
	apiStatusChan := api.APIProcess(ctx, &wg, 5)
	dbStatusChan := db.DBProcess(ctx, &wg, 5)
	redisStatusChan := redis.RedisProcess(ctx, &wg, 5)

	combinedStatus := merge(apiStatusChan, dbStatusChan, redisStatusChan)

	// 채널이 닫힘을 보장 받을 수 있음
	// 종료 채널에 신호를 주는 방식으로 블로킹
	// resultChan := make(chan struct{})
	// go func() {
	// 	wg.Wait()
	// 	resultChan <- struct{}{}
	// }()

	// for {
	// 	select {
	// 	case <-resultChan:
	// 		fmt.Println("작업 완료.")
	// 		return
	// 	case status, ok := <-combinedStatus:
	// 		if !ok {
	// 			continue
	// 		}
	// 		fmt.Println("status 출력 : ", status)
	// 	}
	// }

	for status := range combinedStatus {
		fmt.Println("status 출력 : ", status)
	}
	wg.Wait()
	fmt.Println("작업 종료")
}

func merge(cs ...chan entity.HealthStatus) chan entity.HealthStatus {

	combinedChan := make(chan entity.HealthStatus)

	// 채널에 있는 모든 데이터를 받아 넣음
	var innerWg sync.WaitGroup
	innerWg.Add(len(cs))
	for _, c := range cs {
		go func() {
			defer innerWg.Done()
			for value := range c {
				combinedChan <- value
			}
		}()
	}

	go func() {
		innerWg.Wait()
		fmt.Println("모든 채널 종료")
		close(combinedChan)
	}()

	return combinedChan
}
