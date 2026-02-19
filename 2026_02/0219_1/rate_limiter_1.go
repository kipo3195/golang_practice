package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// 품질 검수 공정 추가
// 품질 검수 장비는 초당 2대 이상 처리하면 과부하가 걸려서 고장이남.
// 생산 -> 타이어 -> 페인팅 -> 검수
// 검수시, 7번째 idx 결함으로 간주, cancel 호출 하고 graceful shutdown
func main() {

	// 생성할 자동차의 수
	makeCarNum := 10

	// 고루틴 관리 wg
	var wg sync.WaitGroup

	// 프로세스 context
	ctx, cancel := context.WithCancel(context.Background())

	// 채널 init
	initTireChan := make(chan *Car)
	paintingChan := make(chan *Car)
	inspectionChan := make(chan *Car, makeCarNum)

	wg.Add(1)
	go makeCar(ctx, &wg, makeCarNum, initTireChan)

	wg.Add(1)
	go initTire(ctx, &wg, initTireChan, paintingChan)

	wg.Add(1)
	go painting(ctx, &wg, paintingChan, inspectionChan)

	wg.Add(1)
	go inspection(cancel, &wg, inspectionChan)

	wg.Wait()

	log.Println("작업완료")
}

func makeCar(ctx context.Context, wg *sync.WaitGroup, makeCarNum int, initTireChan chan<- *Car) {
	defer wg.Done()
	defer close(initTireChan) // 차를 다 만들고 나면 initTire 채널을 닫아버림 -> initTire로직은 채널이 닫혔는지를 체크해야함.

	for i := 1; i <= makeCarNum; i++ {
		// make car 공정
		select {
		case <-ctx.Done():
			log.Println("makeCar 프로세스 종료")
			return
		case initTireChan <- &Car{
			Idx:  i,
			Kind: "suv",
		}:
		}
	}
}

func initTire(ctx context.Context, wg *sync.WaitGroup, initTireChan <-chan *Car, paintingChan chan<- *Car) {
	defer wg.Done()
	defer close(paintingChan)

	for {
		select {
		case <-ctx.Done():
			log.Println("initTire 프로세스 종료")
			return
		case car, ok := <-initTireChan:
			if !ok {
				// 채널이 닫힘
				return
			}
			// tire init 공정
			car.Tire = "hankook"
			select {
			case <-ctx.Done():
				return
			case paintingChan <- car: // 하위 채널이 닫힌 상태
			}
		}
	}

}

func painting(ctx context.Context, wg *sync.WaitGroup, paintingChan <-chan *Car, inspectionChan chan<- *Car) {
	defer wg.Done()
	defer close(inspectionChan)

	for {
		select {
		case <-ctx.Done():
			return
		case car, ok := <-paintingChan:
			if !ok {
				return
			}
			// painting 공정
			car.Color = "red"
			select {
			case <-ctx.Done():
				return
			case inspectionChan <- car: // 검수할 수 있는 상태
			}
		}
	}

}

func inspection(cancel context.CancelFunc, wg *sync.WaitGroup, inspectionChan <-chan *Car) {
	defer wg.Done()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			car, ok := <-inspectionChan
			if !ok {
				return
			}
			if car.Idx == 7 {
				log.Println("결함 발견!!!")
				cancel()
				return
			} else {
				log.Printf("검수 대상 %d번째, Kind :%s, Tire:%s, Color:%s", car.Idx, car.Kind, car.Tire, car.Color)
			}
		}
	}

}

type Car struct {
	Idx   int
	Kind  string
	Tire  string
	Color string
}
