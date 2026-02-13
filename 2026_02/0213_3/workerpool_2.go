package main

import (
	"log"
	"sync"
	"time"
)

// 요구사항
// 자동차 공정 makeCar, initTire, painting
// makeCar, paintng은 0.1초가 걸리지만 initTire의 병목 (1초이상 걸림)
// initTire에 workerpool 도입하여 처리 할 것.

// 핵심 : 내부 공정을 감시하는 고루틴을 별도로 두는 것.

func main() {

	initTireWorkerNum := 3

	// makeCar, painting 종료 감지
	var wg sync.WaitGroup
	// worker pool 종료 감지
	var innerWg sync.WaitGroup

	// car 전달용 채널 생성
	initTireChan := make(chan *Car)
	// paintingChan으로 밀어주는 워커가 3개이므로 painting chan이 동시에 몰릴때 소비하지 못하면 painc
	paintingChan := make(chan *Car, 3)
	//isClosedChan := make(chan bool)

	// 자동차 공정 프로세스
	wg.Add(1)
	go makeCar(&wg, initTireChan /*, isClosedChan*/)
	for i := 1; i <= initTireWorkerNum; i++ {
		// 중간단계의 wg counter를 추가하므로써 모든 고루틴이 종료되었을때 하위 채널을 닫을 수 있도록함.
		innerWg.Add(1)
		go initTire(&innerWg, initTireChan, paintingChan /*, isClosedChan*/)
	}

	// [핵심]
	// inner종료를 감시하는 고루틴을 별도로 빼므로써, main의 병목을 방지함.
	go func() {
		innerWg.Wait()
		// 모든 고루틴이 끝났다 = innerWg의 Done이 모두 호출되므로
		close(paintingChan)
	}()

	wg.Add(1)
	go painting(&wg, paintingChan)

	wg.Wait()
	log.Println("작업종료")

}

func makeCar(wg *sync.WaitGroup, initTireChan chan<- *Car /*, isClosedChan chan<- bool*/) {
	defer wg.Done()
	// 하위 채널 닫아버리기
	defer close(initTireChan)

	// 10대의 차체를 만드는 프로세스
	for i := 1; i <= 10; i++ {
		car := &Car{
			Kind: "suv",
		}
		initTireChan <- car
	}

	// 하나의 worker만 닫도록 -> initTire에서 for루프 deadlock 유발함.
	// isClosedChan <- true

}

// 병목 지점
func initTire(innerWg *sync.WaitGroup, initTireChan <-chan *Car, paintingChan chan<- *Car /*, isClosedChan <-chan bool*/) {
	defer innerWg.Done()

	for car := range initTireChan {
		car.Tire = "kumho"
		time.Sleep(1 * time.Second)
		// 1초가 걸리는 공정
		paintingChan <- car
	}
	// 3개의 고루틴이 모두 닫으려하면 안되니까 -> deadlock 및 paintingChan를 아직 수행하지 못한 고루틴의 panic 유발 (닫힌 채널에 write시)
	// for isclosed := range isClosedChan {
	// 	if isclosed {
	// 		defer close(paintingChan)
	// 	}
	// }
}

func painting(wg *sync.WaitGroup, paintingChan <-chan *Car) {
	defer wg.Done()
	for car := range paintingChan {
		car.Color = "red"

		log.Printf("공정 완료! kind :%s, tire: %s, color:%s \n", car.Kind, car.Tire, car.Color)
	}
}

type Car struct {
	Kind  string
	Tire  string
	Color string
}
