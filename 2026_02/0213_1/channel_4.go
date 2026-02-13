package main

import (
	"log"
	"sync"
	"time"
)

// 0212_3 channel_3.go에 대한 리팩토링
// 종료 수신 채널 대신 close를 이용하여 채널 폐쇄 전파
// for select -> for range
// 매개변수에 수신용 채널인지, 전송용 채널인지 명시적 처리
func main() {

	var wg sync.WaitGroup

	wg.Add(3)

	makeBodyChan := make(chan *Car)
	initTireChan := make(chan *Car)
	paintingChan := make(chan *Car)

	go makeCar(&wg, makeBodyChan, initTireChan)
	go initTire(&wg, initTireChan, paintingChan)
	go painting(&wg, paintingChan)

	// 생산자 소비자 패턴에서 데이터를 생산하는 역할.
	for i := 1; i <= 10; i++ {
		if i%2 == 1 {
			makeBodyChan <- &Car{}
		} else {
			makeBodyChan <- &Car{}
		}
		time.Sleep(1 * time.Second)
	}

	// 채널 닫아버리기 - 하위 채널의 닫기 처리를 위함
	close(makeBodyChan)

	wg.Wait()
	log.Println("작업 완료 ")
}

func makeCar(wg *sync.WaitGroup, makeBodyChan <-chan *Car, initTireChan chan<- *Car) {
	log.Println("작업 시작")
	defer wg.Done()

	// 해당 함수 종료시 하위 채널 닫아버리기
	defer close(initTireChan)

	// for select 에서 for range로 변경한 이유
	// for select : 해당 채널에 데이터가 들어올때까지 기다림므로 for 루프 종료 X
	// for range : 해당 채널에 데이터가 들어오지 않는다면 for 루프 종료 (채널이 닫혔다 = 데이터 X)이니까 루프종료
	for car := range makeBodyChan {
		now := time.Now().Second()
		devide := now % 2
		var color string
		if devide == 0 {
			car.Kind = "suv"
		} else {
			car.Kind = "sport"
		}
		car.Color = color
		initTireChan <- car
	}
}

func initTire(wg *sync.WaitGroup, initTireChan <-chan *Car, paintingChan chan<- *Car) {
	defer wg.Done()
	defer close(paintingChan)

	for car := range initTireChan {
		now := time.Now().Second()
		devide := now % 2
		if devide == 0 {
			car.Tire = "normal"
		} else {
			car.Tire = "snow"
		}
		paintingChan <- car
	}
}

func painting(wg *sync.WaitGroup, paintingChan <-chan *Car) {
	defer wg.Done()
	for car := range paintingChan {
		now := time.Now().Second()
		devide := now % 2
		if devide == 0 {
			car.Color = "red"
		} else {
			car.Color = "blue"
		}
		log.Printf("make end. kind :%s, color:%s, tire :%s \n", car.Kind, car.Color, car.Tire)
	}
}

// 자동차 구조체
type Car struct {
	Kind  string
	Tire  string
	Color string
}
