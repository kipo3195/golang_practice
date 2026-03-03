package process

import (
	"sync"
)

type Process struct {
	workerNum  int
	ResultChan chan int
	initChan   chan int
}

func NewProcess(workerNum int, initChan chan int) *Process {

	resultChan := make(chan int, 20)
	return &Process{
		workerNum:  workerNum,
		ResultChan: resultChan,
		initChan:   initChan,
	}
}

func (r *Process) Work() {

	// ticker := time.NewTicker(10 * time.Millisecond)
	// defer ticker.Stop()

	var wg sync.WaitGroup
	for i := 1; i <= r.workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				//case <-ticker.C:
				// [TODO] generator에서 initChan에 데이터를 모두 넣고 닫았으니
				// initChan은 for-range여도 된다 (ok체크 할 필요 없음 )
				case v, ok := <-r.initChan: // 채널에서 데이터를 뽑아내자
					// 단, ok는 채널이 닫혀서 정상이 아님
					if !ok {
						// 닫힌 채널에서 데이터를 계속 뽑아내면 nil임. 그
						return
					}
					sqare := v * v
					r.ResultChan <- sqare // resultChan의 cap은 20
					// [오답 체크 - 1]
					// 데이터의 개수가 20보다 커지게 되면 ResultChan에 데이터를 넣다가 멈춤
					// 아래의 wg.Wait()이 풀리지 않음.
					// ResultChan이 소비되지 않으므로..
					// 딱 개수에 맞게만 처리가능함.
				}
			}
		}()
	}

	// [TODO] 동기적으로 처리되지 않으려면 어떻게 해야할지 생각해보자
	wg.Wait()
	// 모든 연산이 끝났으므로 initChan 닫기
	close(r.ResultChan)
}
