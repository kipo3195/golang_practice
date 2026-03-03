package worker

import (
	"context"
	"log"
	"sync"
	"test/job"
	"time"
)

func Worker(ctx context.Context /*cancel context.CancelFunc,*/, wg *sync.WaitGroup, workChan <-chan job.Job) {
	defer wg.Done()

	// 결국 들어오는 데이터는 정해져있기 때문에 for-select로
	// 닫히는 것이 명확하다면 for range가 더나음 어차피 소비하는건 동일하므로
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// 채널의 데이터를 모두 뺄때까지
	for value := range workChan {
		select {
		case <-ctx.Done():
			// 타임아웃인 경우
			log.Println("작업종료!")
			// cancel()
			// ctx.Done()이 이미 호출되었다는 건,
			// 상위(Main)에서 이미 취소나 타임아웃이 발생했다는 뜻입니다. 자식이 부모에게 "취소해줘"라고 다시 말할 필요는 없습니다.
			return
		case <-ticker.C:
			// 티커의 동작
			log.Printf("new :%d", value)
			// select {
			// case <-ctx.Done():
			// case news := <-workChan:
			// 	log.Printf("news :%d", news.Id)
			// }
		}
	}

}
