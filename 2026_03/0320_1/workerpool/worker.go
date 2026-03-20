package workerpool

import (
	"context"
	"fmt"
	"sync"
	"test/entity"
	"time"
)

type WorkerPool struct {
	imageChan  chan entity.Job
	ResultChan chan entity.Result
	workerNum  int
}

func NewWorkerPool(imageChan chan entity.Job, workerNum int) *WorkerPool {

	resultChan := make(chan entity.Result, workerNum)
	return &WorkerPool{
		imageChan:  imageChan,
		ResultChan: resultChan,
		workerNum:  workerNum,
	}
}

func (r *WorkerPool) Worker(ctx context.Context, wg *sync.WaitGroup) {

	wg.Add(r.workerNum)
	for i := 1; i <= r.workerNum; i++ {

		go func(idx int) {
			defer wg.Done()
			for {
				select {
				case img, ok := <-r.imageChan:
					if !ok {
						fmt.Printf("%d worker process end. all process end\n", idx)
						return
					}
					workingTime := img.Complexity * 200
					time.Sleep(time.Duration(workingTime) * time.Millisecond)

					// sleep에서 깨어나서 처리 할 수 있는가
					select {
					case <-ctx.Done():
						fmt.Printf("%d worker process end. timeout\n", idx)
						return
					case r.ResultChan <- entity.Result{
						ID:          img.ID,
						WorkingTime: workingTime,
						ThreadID:    idx,
					}:

					}

				case <-ctx.Done():
					fmt.Printf("%d worker process end. timeout\n", idx)
					return
				}
			}

		}(i)
	}

}
