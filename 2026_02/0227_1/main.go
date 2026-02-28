package main

import (
	"chat/chat"
	"chat/entity"
	"context"
	"fmt"
	"log"
)

// 메시지를 비동기로 처리하는 서버를 구현하라.
// https://chatgpt.com/c/69a19bcf-49cc-8322-8911-37c71807ee2d

func main() {
	// handler 생성
	handler := chat.NewChatHandler()

	// context 생성
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 채팅 처리 프로세스 생성
	p := chat.NewChatProcessor(3, handler, cancel)

	p.Start(context)

	for i := 1; i <= 10; i++ {

		// payload를 초기화 해 주지않으면 p.Submit을 수행하는건 고루틴이므로 클로저 처리됨 -> 나중에 수행되므로...
		payload := fmt.Sprintf("안녕하세요 %d", i)
		log.Println("payload : ", payload)
		msg := entity.Message{
			ID:      "1",
			Payload: payload,
		}
		p.Submit(msg)
	}

	// time.Sleep(5 * time.Second)
	p.Stop()

	log.Println("작업 종료.")
}
