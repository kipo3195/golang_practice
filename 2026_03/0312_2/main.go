package main

import (
	"context"
	"fmt"
	"sync"
	"test/flusher"
	"test/producer"
	"time"
)

// 왜 빈 메시지가 출력되는지 점검
// https://gemini.google.com/share/f5a713edc95a
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	// 20개의 이벤트 주입
	eventChan := producer.Produce(ctx, 20)
	wg.Add(1)
	go func() {
		defer wg.Done()
		flusher.Flusher(ctx, eventChan)
	}()

	wg.Wait()

	fmt.Println("작업 종료")

}
