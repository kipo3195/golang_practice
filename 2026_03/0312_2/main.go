package main

import (
	"context"
	"test/flusher"
	"test/producer"
	"time"
)

// 왜 빈 메시지가 출력되는지 점검
// https://gemini.google.com/share/f5a713edc95a
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 20개의 이벤트 주입
	eventChan := producer.Produce(ctx, 20)
	flusher.Flusher(ctx, eventChan)

}
