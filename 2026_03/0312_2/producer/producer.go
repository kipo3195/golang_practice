package producer

import (
	"context"
	"fmt"
	"test/entity"
)

func Produce(ctx context.Context, eventNum int) chan entity.Message {

	eventChan := make(chan entity.Message)
	go func(eventNum int) {
		defer close(eventChan)
		for i := 1; i <= eventNum; i++ {
			select {
			case <-ctx.Done():
				// 채널 닫기
				return
			case eventChan <- entity.Message{
				Payload: fmt.Sprintf("이벤트 %d", i),
			}:
			}
		}
	}(eventNum)

	return eventChan
}
