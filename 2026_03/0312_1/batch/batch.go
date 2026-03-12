package batch

import (
	"context"
	"fmt"
	"test/entity"
)

type Batch struct {
	recvChan  chan entity.Message
	batchChan chan entity.FlushEntity
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
					r.BatchProcess(ctx, msg)
				}
			}
		}

	}()

}

// batch -> collect or flush
func (r *Batch) BatchProcess(ctx context.Context, msg entity.Message) {

	select {

	case b := <-r.batchChan:

		if b.IsFlush {
			// DB 프로세스 실행

			dbSaveProcess(b.Messages)

			// 새로운 Flush entity를 생성하거나 isFlush, Messages를 비어있는 걸로 변경한 후 ---------- 2
			// 아래의 1번 로직 수행
		}

		b.Messages = append(b.Messages, &msg)
		// batch chan에 넣고나서 상태를 확인, 5개면 is flush로 변경 -------- 1 (반복)
	case <-ctx.Done():
		// 여기서는 남아있는 모든 작업을 수행 할 수 있는 함수를 생성하여 처리하는것이 어떨까?
		// 아니면 RecvProcess에서 모든 작업을 수행 할 수 있는 함수를 호출하는게 맞는건가
		// 근데 그러면 너무 많은 로직에서 처리를 해야될거같은데 r.batchChan을 닫지 않고 정리하는 한 다음 채널을 닫는 함수를 생성?
		return
	}

}

func dbSaveProcess(msgs []*entity.Message) {

	for i := 1; i <= len(msgs); i++ {
		fmt.Printf("DB 저장처리 로직 수행 : %s \n", msgs[i].Payload)
	}

}
