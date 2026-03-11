package workerpool

import (
	"context"
	"sync"
	"test/entity"
)

type WorkerPool struct {
	Workers []*Worker
}

func NewWorkerPool(workerNum int) *WorkerPool {

	workers := make([]*Worker, 0)

	for i := 1; i <= workerNum; i++ {

		w := &Worker{
			// 워커의 로컬 맵
			pendingMap: make(map[string]entity.Task),
			EventChan:  make(chan entity.Task),
		}
		workers = append(workers, w)
	}

	return &WorkerPool{
		Workers: workers,
	}
}

func (r *WorkerPool) Init(ctx context.Context, wg *sync.WaitGroup) {

	for idx, w := range r.Workers {
		wg.Add(1)
		// 작업 처리 워커
		go w.Work(ctx, wg, idx)
	}

}
