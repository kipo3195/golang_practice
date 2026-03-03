package generator

import "log"

type Generator struct {
	n int
}

func NewGenerator(n int) chan int {

	initChan := make(chan int, n)

	// [TODO] 여기서 채널을 비동기 처리하고 채널을 닫더라도
	// 하위로직에서 initChan의 데이터가 모두 소진될 수 있다
	for i := 1; i <= n; i++ {
		initChan <- i
	}

	close(initChan)
	log.Println("NewGenerator end")
	return initChan
}
