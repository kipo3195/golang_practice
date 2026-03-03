package worker

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

func Worker(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, resultChan chan<- string, a string) {
	defer wg.Done()

	min := 100
	max := 500
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond
	log.Printf("작업 : %s 예상시간 : %s\n", a, sleepDuration)
	select {
	case <-ctx.Done():
		// 이미 다른 누군가가 실패해서 취소됨 OR 타임아웃
		// 취소요청 들어감
		return
	case <-time.After(sleepDuration):
		// 여기서 시간이 지났는지 점검해야할거같은데 X
		// context init시에 잡은 시간이 지나면 자동으로 ctx.Done()호출함.

		// 문제는 특정 작업은 무조건 실패!
		if a == "lierCheck" {
			cancel()
			return
		}

		select {
		case <-ctx.Done():
			// 취소요청 들어감
			return
		case resultChan <- a:
		}

	}

}
