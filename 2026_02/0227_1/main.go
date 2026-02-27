package main

import (
	"chat/chat"
	"chat/entity"
	"context"
	"log"
)

// 메시지를 비동기로 처리하는 서버를 구현하라.

func main() {
	// handler 생성
	handler := chat.NewChatHandler()

	// 채팅 처리 프로세스 생성
	p := chat.NewChatProcessor(3, handler)

	// context 생성
	context := context.Background()

	p.Start(context)

	for i := 1; i <= 10; i++ {
		msg := entity.Message{
			ID:      "1",
			Payload: "안녕하세요",
		}
		p.Submit(msg)
	}

	p.Stop()

	log.Println("작업 종료.")
}
