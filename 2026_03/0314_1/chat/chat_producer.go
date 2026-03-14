package chat

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ChatProduceProcess(ctx context.Context, wg *sync.WaitGroup, processNum int) chan string {

	chatChan := make(chan string)
	ticker := time.NewTicker(2 * time.Second)

	go func(processNum int) {
		defer wg.Done()
		defer close(chatChan)
		defer ticker.Stop()

		for i := 1; i <= processNum; i++ {
			// 설정된 시간 만큼 대기
			<-ticker.C
			fmt.Println("[Chat] logging")
			chatChan <- fmt.Sprintf("chat log %d", i)
		}

	}(processNum)

	return chatChan
}
