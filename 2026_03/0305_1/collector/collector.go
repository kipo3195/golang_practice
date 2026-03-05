package collector

import (
	"context"
	"log"
	"test/entity"
)

type Collector struct {
	msgChan chan<- entity.Message
}

func NewCollector(msgChan chan<- entity.Message) *Collector {
	return &Collector{
		msgChan: msgChan,
	}
}

func (r *Collector) Collect(ctx context.Context, msg []string) {

	// 클로저
	go func(target []string) {
		// 모든 메시지 전송후 닫기
		defer close(r.msgChan)

		for idx, value := range target {
			select {
			case <-ctx.Done():
				// 외부 에러로 인한 로직 종료
				return
			case r.msgChan <- entity.Message{
				ID:      idx,
				Content: value,
			}:
				log.Println("value :", value)
			}
		}
	}(msg)

}
