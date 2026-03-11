package workerpool

import (
	"context"
	"fmt"
	"sync"
	"test/entity"
)

type Worker struct {
	// 자신의 로컬 맵
	pendingMap map[string]entity.Task
	// 자신의 이벤트 채널
	EventChan chan entity.Task
}

func (r *Worker) Work(ctx context.Context, wg *sync.WaitGroup, idx int) {
	defer wg.Done()

	// 작업 종료시 워커의 채널 닫기
	defer close(r.EventChan)

	for {
		select {
		case <-ctx.Done():
			// 외부 호출에 의한 이벤트 종료
			return
		case task := <-r.EventChan:
			fmt.Printf("%d 워커가 처리하는 task :%s", idx, task.Payload)
		}
	}

}
