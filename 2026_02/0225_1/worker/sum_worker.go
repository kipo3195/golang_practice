package worker

type SumWorker struct {
	recvChan chan int
}

func NewSumWorker() SumWorker {
	return SumWorker{}
}
