package batch

import (
	"context"
	"fmt"
	"test/entity"
	"time"
)

type Batch struct {
	recvChan   chan entity.Message
	batchChan  chan entity.FlushEntity
	pendingMap map[string]*entity.SaveEntity
}

func NewBatch() *Batch {

	return &Batch{
		recvChan:  make(chan entity.Message),
		batchChan: make(chan entity.FlushEntity),
	}
}

// 수신 -> batch
func (r *Batch) RecvProcess(ctx context.Context, msgChan <-chan entity.Message) {

	go func() {

		defer close(r.recvChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("외부에 의한 종료 호출. Batch Process 종료")
				// 여기에 모든 배치 작업을 마무리 하는 로직이 들어가야되나?
				// 근데 그러면 너무 많은 로직에서 처리를 해야될거같은데 r.batchChan을 닫지 않고 정리하는 한 다음 채널을 닫는 함수를 생성?
				return
			case msg := <-msgChan:
				select {
				case <-ctx.Done():
					// 취소 신호
					// 여기에 모든 배치 작업을 마무리 하는 로직이 들어가야되나?
					// 근데 그러면 너무 많은 로직에서 처리를 해야될거같은데 r.batchChan을 닫지 않고 정리하는 한 다음 채널을 닫는 함수를 생성?
				default:
					// 다시 job에 넣어버릴까?
					r.batchChan <- entity.FlushEntity{
						IsFlush:  false,
						Messages: &msg,
					}

				}
			}
		}

	}()

}

// batch -> collect or flush
func (r *Batch) BatchProcess(ctx context.Context, msg entity.Message) {

	go func() {

		for job := range r.batchChan {

			if job.IsFlush {
				savedEntity := r.pendingMap["t"]
				dbSaveProcess(savedEntity.Msg)
				// 전송 했으니 삭제
				delete(r.pendingMap, "t")
				return
			}

			r.addProcess(ctx, job.Messages)
		}

	}()

}

func (r *Batch) addProcess(ctx context.Context, msg *entity.Message) {

	v, exists := r.pendingMap["t"]

	if !exists {
		// 최초 혹은 발송 후
		temp := make([]*entity.Message, 0)
		temp = append(temp, msg)

		saveEntity := &entity.SaveEntity{
			Msg: temp,
		}

		saveEntity.Timer = time.AfterFunc(3*time.Second, func() {
			// 중요: 워커 자신의 채널로 Flush 작업을 다시 던집니다.
			r.batchChan <- entity.FlushEntity{
				IsFlush:  true,
				Messages: nil,
			}
		})

		return
	}

	v.Msg = append(v.Msg, msg)
}

func dbSaveProcess(msgs []*entity.Message) {

	for i := 1; i <= len(msgs); i++ {
		fmt.Printf("DB 저장처리 로직 수행 : %s \n", msgs[i].Payload)
	}

}
