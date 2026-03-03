package worker

import (
	"context"
	"math/rand"
	"sync"
	"test/dto"
	"time"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, jobName string, resultChan chan<- dto.Result) {
	defer wg.Done() // 함수 종료 시 무조건 실행 보장, for 루프에서는 return 처리만.
	// 만약 로직이 복잡해져서 다른 곳에서 return 해버리면 Done()이 호출되지 않아 메인이 영원히 기다리게 됩니다.

	// time.Sleep은 Context 취소를 감지하지 못합니다.
	// 지정된 시간 동안 절대 깨어나지 않습니다. 만약 타임아웃이 500ms이고 슬립이 800ms라면, 500ms 시점에 취소 신호가 와도 무시하고 800ms까지 푹 자버립니다.
	// select와 time.After를 써야 즉각 반응할 수 있습니다.

	min := 100
	max := 800
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond

	select {
	case <-ctx.Done():
		return
	case <-time.After(sleepDuration): // 시간이 지난 후 실행 (time.Sleep은 취소시간 감지 불가!!)
		var err string
		if jobName == "fail" {
			err = "error"
		}

		// 전송 시점에도 Context 확인 (보내는 도중 취소될 수 있으므로)
		select {
		case <-ctx.Done():
			return
		case resultChan <- dto.Result{
			Value: jobName,
			Err:   err,
		}:
		}
	}

}
