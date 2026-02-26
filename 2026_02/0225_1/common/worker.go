package common

import (
	"sync"
)

type Worker struct {
	recvChan chan int
	op       Operator
}

func NewWorker(recvChan chan int, op Operator) *Worker {
	return &Worker{
		recvChan: recvChan,
		op:       op,
	}
}

func (r *Worker) Process(wg *sync.WaitGroup, resultChan chan<- int) {
	defer wg.Done()

	for data := range r.recvChan {
		resultChan <- r.op(data)
	}
}
