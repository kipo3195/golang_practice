package collector

type Collector struct {
	resultChan chan int
	MergeChan  chan int
}

func NewCollector(resultChan chan int, chanSize int) *Collector {
	mergeChan := make(chan int, 20)
	return &Collector{
		resultChan: resultChan,
		MergeChan:  mergeChan,
	}
}

// [TODO] // 채널 데이터를 받아서 슬라이스로 모아주는 역할
func (r *Collector) Merge() chan int {

	for value := range r.resultChan {
		r.MergeChan <- value
	}
	close(r.MergeChan)
	return r.MergeChan
}
