package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"test/entity"
	"test/external"
	"time"
)

type Dispatcher struct {
	workerNum  int
	msgChan    <-chan entity.Message // 메시지 수신용 채널
	resultChan chan<- entity.Result  // 결과 합계 전송용 채널
}

func NewDispatcher(workerNum int, msgChan <-chan entity.Message, resultChan chan<- entity.Result) *Dispatcher {
	return &Dispatcher{
		workerNum:  workerNum,
		msgChan:    msgChan,
		resultChan: resultChan,
	}
}

func (r *Dispatcher) Work(ctx context.Context) {

	// msg 슬라이스
	var wg sync.WaitGroup
	for i := 1; i <= r.workerNum; i++ {
		wg.Add(1)
		go func(workerId int) {
			// wg.Done()을 반복해서 사용하는 습관 지양하기. 전역적으로 한곳에서만 처리 할 수 있도록
			defer wg.Done()
			r.workerLoop(ctx, workerId)
		}(i)
	}

	// 모든 고루틴이 종료되길 대기
	go func() {
		wg.Wait()
		close(r.resultChan)
	}()

}

func (r *Dispatcher) workerLoop(ctx context.Context, workerID int) {

	for {
		select {
		// 외부 이벤트에 의한 종료 점검용
		case <-ctx.Done():
			return
		case msg, ok := <-r.msgChan: // 종료되지 아니한 상태
			if !ok {
				// 채널이 닫혔는지 여부를 판단.
				return
			}
			// 별도 ctx 생성 (API의 타임아웃 감지)
			// 무조건 새로운 ctx를 만드는게 능사가 아님. 기존 ctx와 연관이 있어야할때는 부모 wrapping하여 생성.
			// ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
			// [수정 3] 부모 context(ctx)를 이어받아야 Graceful Shutdown이 됨
			taskCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)

			result := entity.Result{Success: true}

			err := external.SendPushNotification(taskCtx, msg)
			if err != nil {
				result.Success = false
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("worker %d 타임아웃에 의한 실패 msg : %q", workerID, msg)
				}
			}
			// defer로 cancel하게되면 workerLoop함수가 끝난뒤에 taskCtx 가 종료됨 -> cancel 함수가 스택에 쌓임.
			cancel()

			select {
			// 외부 이벤트에 종료되어야 하는지 점검
			case <-ctx.Done():
				return
			case r.resultChan <- result:
			}
		}
	}
}
