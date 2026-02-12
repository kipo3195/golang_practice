package main

import (
	"log"
	"sync"
	"time"
)

// 로깅 데이터 수신용 채널 + 로깅 데이터 수신용 채널 닫기용 채널 생성
func main() {

	var wg sync.WaitGroup

	dataChan := make(chan int)
	endChan := make(chan struct{})

	wg.Add(1)
	go print(&wg, dataChan, endChan)

	log.Println("1초 간격으로 dataChan에 데이터 전송 ")
	for i := 0; i <= 5; i++ {
		dataChan <- i
		time.Sleep(1 * time.Second)
	}

	log.Println("5초 동안 슬립 후 endChan에 데이터 추가 ")
	time.Sleep(5 * time.Second)
	endChan <- struct{}{}

	wg.Wait() // 여기서 기다려라
	log.Println("작업 종료")
}

func print(wg *sync.WaitGroup, dataChan chan int, endChan chan struct{}) {

	for {

		select {
		case i := <-dataChan:
			log.Println("data 채널에 데이터 수신 :", i)
		case <-endChan: // 신호 자체가 중요할때는 변수에 할당하지 않고 <- endChan만으로 처리함.
			log.Println("end 채널에 데이터 수신됨.")
			wg.Done()
			return // 필수!!!! return하지 않는다면 wg.Done() 호출 후 다시 select 대기 해버림..
		}
	}

}
