package producer

import (
	"context"
	"fmt"
	"test/entity"
)

func Produce(ctx context.Context, messageNum int) chan<- entity.Message {

	msgChan := make(chan entity.Message)

	go func(messsageNum int) {
		// 다 쓴다음 채널 닫기
		defer close(msgChan)

		for i := 1; i <= messageNum; i++ {

			select {
			case msgChan <- entity.Message{
				Payload: fmt.Sprintf("%d 번째 메시지 입니다.", i),
			}:
			case <-ctx.Done():
				fmt.Println("외부에 의한 종료 호출. Produce 종료")
				return
			}
		}

	}(messageNum)

	return msgChan
}
