package flusher

import (
	"context"
	"fmt"
	"test/entity"
	"time"
)

func Flusher(ctx context.Context, eventChan <-chan entity.Message) {

	// 메시지가 5개 쌓이면 즉시 저장합니다.
	const batchSize = 5
	var buffer []entity.Message

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 메시지 저장용

	for {
		select {
		case <-ctx.Done():
			// ctx로 종료되거나, 잔여 이벤트 모두 DBSave
			fmt.Println("컨텍스트 만료: 남은 데이터를 플러시하고 종료합니다.")
			buffer = DBSave(buffer)
			return
		case msg, ok := <-eventChan:
			if !ok {
				// 남은 이벤트 모두 저장
				if len(buffer) > 0 {
					buffer = DBSave(buffer)
					return
				}
			}
			// 채널에 데이터 수신, 5개 이상이라면 DBSave
			buffer = append(buffer, msg)
			if len(buffer) == batchSize {
				// flush 해야되거나
				fmt.Print("[Size Limit] ")
				buffer = DBSave(buffer)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				fmt.Print("[Time Limit] ")
				buffer = DBSave(buffer)
			}
		}
	}

}

func DBSave(buffer []entity.Message) []entity.Message {
	if len(buffer) == 0 {
		return buffer
	}
	fmt.Printf("DB 저장 완료: %d개의 메시지 (%v)\n", len(buffer), buffer)
	// 빈배열 반환
	return []entity.Message{}
}
