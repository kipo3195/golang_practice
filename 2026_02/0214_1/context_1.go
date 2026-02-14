package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// 자동차 공정 프로세스
// makeCar -> initTire -> painting
// painting 5번째 차량에서 화재 발생, makeCar, initTire는 채널 데이터 소모하지 않고 즉시 중단하기

// [핵심]
// **"고루틴이 채널 때문에 대기(Blocking)할 가능성이 있는 모든 지점에 탈출구를 만든다"**는 전략입니다.

func main() {

	// 사고 발생시 처리를 위한 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var innerWg sync.WaitGroup

	initTireWorkerNum := 3

	initTireChan := make(chan *Car)
	paintingChan := make(chan *Car, 3)

	wg.Add(1)
	go makeCar(ctx, &wg, initTireChan)

	for i := 1; i <= initTireWorkerNum; i++ {
		wg.Add(1)
		innerWg.Add(1)
		go initTire(ctx, &wg, &innerWg, initTireChan, paintingChan)
	}

	wg.Add(1)
	go painting(ctx, cancel, &wg, paintingChan)

	go func() {
		innerWg.Wait()
		// innnerWg가 모두 종료시 하위 채널을 닫음
		close(paintingChan)
	}()

	wg.Wait()
	log.Println("작업 종료.")
}

// 전송용 채널이라는 것을 명시
func makeCar(ctx context.Context, wg *sync.WaitGroup, initTireChan chan<- *Car) {

	defer wg.Done()
	defer close(initTireChan)

	for i := 1; i <= 10; i++ {

		select {
		case <-ctx.Done(): // **"채널이 닫혔음(Close)"**을 수신하는 것입니다.
			//context.WithCancel(parent)을 호출하면 내부적으로 Done이라는 이름의 채널이 생성됩니다.
			//어디선가 cancel() 함수를 실행하면, Go 런타임은 이 Done 채널을 close(ch) 해버립니다.
			// Go 채널의 특징 중 하나는 **"닫힌 채널로부터의 수신은 즉시 완료된다"**는 것입니다.
			log.Println("[makeCar] 화재로 인한 공정 종료")
			return

		case initTireChan <- &Car{
			Kind: "suv",
			Idx:  i,
		}:
		}
	}

}

// 수신용 채널이라는 것을 명시, 전송용 채널이라는 것을 명시
func initTire(ctx context.Context, wg *sync.WaitGroup, innerWg *sync.WaitGroup, initTireChan <-chan *Car, paintingChan chan<- *Car) {
	defer innerWg.Done()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// 취소 신호가 왔을때
			log.Println("[initTire] 화재 발생으로 인한 작업 중지 ")
			return
		case car, ok := <-initTireChan:
			if !ok {
				log.Println("채널이 닫힘.")
				return
			}
			car.Tire = "hankook"
			time.Sleep(1 * time.Second)

			select {
			case paintingChan <- car: // 채널에 데이터를 넣을수 있는 상태라면
			case <-ctx.Done(): // 화재 발생이라면 -- painting함수에서 cancel() 호출된 상태
				log.Println("[initTire] 데이터 전송 중 화재 발생 - 중단")
				return
			}
		}
	}

}

func painting(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, paintingChan <-chan *Car) {
	defer wg.Done()

	for car := range paintingChan {
		if car.Idx == 5 {
			log.Println("화재 발생!!!!!!!!!!!!!!!!!!")
			cancel() // <-ctx.Done()에 이벤트 할당
			return
		}
		car.Color = "red"
		log.Printf("차량번호 idx %d, kind :%s, tire :%s, color:%s", car.Idx, car.Kind, car.Tire, car.Color)
	}

}

type Car struct {
	Kind  string
	Tire  string
	Color string
	Idx   int
}
