package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	workerNum int
}

func NewWorkerPool(workerNum int) *WorkerPool {
	return &WorkerPool{
		workerNum: workerNum,
	}
}

func (r *WorkerPool) Work(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan int, resultChan chan<- int) {

	for i := 1; i <= r.workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					fmt.Println("타임 아웃 혹은 취소")
					return
				case data, ok := <-dataChan:
					if !ok {
						// 모든 데이터를 다 뽑아냈다
						fmt.Println("모든 데이터를 다 뽑아낸거 같네요.")
						return
					}
					time.Sleep(100 * time.Millisecond)

					select {
					case resultChan <- (data * 2):
					case <-ctx.Done():
						fmt.Println("타임 아웃 혹은 취소")
						return
					}
				}
			}
		}()
	}

}
