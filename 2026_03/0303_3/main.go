package main

import (
	"context"
	"log"
	"sync"
	"test/job"
	"test/worker"
	"time"
)

// 당신은 회원 20명에게 뉴스레터를 보내야 합니다. 하지만 메일 발송 서버는 과부하를 막기 위해 **"0.5초에 1건씩"**만 요청을 받습니다.
// Job Queue: 20개의 메일 발송 작업(ID: 1~20)을 채널에 넣습니다.
// Rate Limit: 워커는 작업을 처리하기 전에 **반드시 0.5초(500ms)의 대기 시간(Tick)**을 가져야 합니다. (time.Sleep 대신 time.Ticker 사용 권장)
// Timeout: 전체 작업은 5초 안에 끝나야 합니다. (20개를 0.5초씩 보내면 10초가 걸리므로, 중간에 타임아웃으로 끊겨야 정상입니다.)
// Graceful Shutdown: 타임아웃 신호(ctx.Done)가 오면, 대기 중이던 작업도 즉시 멈추고 종료해야 합니다.

func main() {

	// job 채널에 데이터 넣으면
	// worker가 job 채널에서 소비.
	// 처리로직에서 Ticker 사용

	// 전체 작업 타이머
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 작업
	workChan := make(chan job.Job, 20)

	var wg sync.WaitGroup

	// worker
	wg.Add(1)
	go worker.Worker(ctx, &wg, workChan)

	for i := 1; i <= 20; i++ {
		workChan <- job.Job{
			Id: i,
		}
	}

	// 데이터를 다 넣고 닫아버림 -> worker에 for-range 이므로
	// 닫힐때까지 처리함.
	// for select로 했더니 타임아웃 시간을 증가시켰을때
	// ctx.Done()은 호출 되지 않고 계속 ticker는 울림
	// for-range로 채널이 닫힐때까지라면 시간이 증가해도 채널에 있는 데이터를 다 빼면 끝나므로 효율적임.

	close(workChan)

	// 모든 고루틴이 종료된다면, 채널을 닫도록함
	wg.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		log.Println(" 시간초과로 인한 작업 종료 ")
	} else {
		log.Println("작업 끝")
	}

}
