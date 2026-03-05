package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	go func() {

		for i := 1; i <= r.workerNum; i++ {
			wg.Add(1)
			for {
				select {
				// 외부 이벤트에 의한 종료
				case <-ctx.Done():
					wg.Done()
					return
				case msg, ok := <-r.msgChan: // 종료되지 아니한 상태
					log.Println("dispacher 내부 : ", msg, ok)
					if !ok {
						wg.Done()
						return
					}
					// 별도 ctx 생성 (API의 타임아웃 감지)
					ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
					defer cancel2()
					select {
					case <-ctx.Done():
						wg.Done()
						return
					default:
						err := external.SendPushNotification(ctx2, msg)
						if err != nil {
							// result entity에 실패 ++
							r.resultChan <- entity.Result{
								Success: false,
							}
							if errors.Is(err, context.DeadlineExceeded) {
								fmt.Println("타임아웃에 의한 실패 msg : ", msg)
							}
						}
						// result entity에 성공 ++
						r.resultChan <- entity.Result{
							Success: true,
						}
					}
				}
			}
		}
	}()

	// 모든 고루틴이 종료되길 대기
	go func() {
		wg.Wait()
		close(r.resultChan)
	}()

}
