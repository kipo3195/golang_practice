package main

import (
	"log"
	"sync"
	"time"
)

// 자동차 생산 라인 가동 채널
// 프로세스 차체 생산 -> 바퀴 조립 -> 페인팅을 채널을 이용하여 처리
// 생산자 - 소비자 패턴 : mutex를 사용하지 않고 채널을 이용하여 race condition을 방지하는 방법
func main() {

	var wg sync.WaitGroup

	makeBodyChan := make(chan *Car)
	makeEnd := make(chan struct{})
	initTireChan := make(chan *Car)
	initEnd := make(chan struct{})
	paintingChan := make(chan *Car)
	paintingEnd := make(chan struct{})

	wg.Add(3) // 3개의 고루틴 실행 - Add < Done() : 음수가 될 수 없으므로 panic, Add > Done() DeadLock
	go makeCar(&wg, makeEnd, makeBodyChan, initEnd, initTireChan)
	go initTire(&wg, initEnd, initTireChan, paintingEnd, paintingChan)
	go painting(&wg, paintingEnd, paintingChan)

	// 생산자 소비자 패턴에서 데이터를 생산하는 역할.
	for i := 1; i <= 10; i++ {
		if i%2 == 1 {
			makeBodyChan <- &Car{}
		} else {
			makeBodyChan <- &Car{}
		}
		time.Sleep(1 * time.Second)
	}

	makeEnd <- struct{}{}

	wg.Wait()
	log.Println("작업 완료 ")
}

func makeCar(wg *sync.WaitGroup, makeEnd chan struct{}, makeBodyChan chan *Car, initEnd chan struct{}, initTireChan chan *Car) {
	log.Println("작업 시작 ")
	for {
		select {
		case car := <-makeBodyChan:
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
		case <-makeEnd:
			wg.Done()
			close(makeBodyChan)
			initEnd <- struct{}{}
			return
		}
	}

}

func initTire(wg *sync.WaitGroup, initEnd chan struct{}, initTireChan chan *Car, paintingEnd chan struct{}, paintingChan chan *Car) {
	for {
		select {
		case car := <-initTireChan:
			now := time.Now().Second()
			devide := now % 2
			if devide == 0 {
				car.Tire = "normal"
			} else {
				car.Tire = "snow"
			}
			paintingChan <- car
		case <-initEnd:
			wg.Done()
			close(initTireChan)
			paintingEnd <- struct{}{}
			return
		}
	}
}

func painting(wg *sync.WaitGroup, paintingEnd chan struct{}, paintingChan chan *Car) {
	for {
		select {
		case car := <-paintingChan:
			now := time.Now().Second()
			devide := now % 2
			if devide == 0 {
				car.Color = "red"
			} else {
				car.Color = "blue"
			}
			log.Printf("make end. kind :%s, color:%s, tire :%s \n", car.Kind, car.Color, car.Tire)
		case <-paintingEnd:
			wg.Done()
			close(paintingChan)
			return
		}
	}
}

// 자동차 구조체
type Car struct {
	Kind  string
	Tire  string
	Color string
}
