package chat

import (
	"chat/entity"
	"chat/handler"
	"chat/processor"
	"chat/worker"
	"context"
	"errors"
	"sync"
)

type ChatProcessor struct {
	workerPoolNum int
	msgChan       chan entity.Message
	worker        []*worker.ChatWorker
	wg            *sync.WaitGroup
	handler       handler.Handler
}

func NewChatProcessor(workerPoolNum int, handler handler.Handler) processor.Processor {
	var wg sync.WaitGroup
	return &ChatProcessor{
		workerPoolNum: workerPoolNum,
		msgChan:       make(chan entity.Message),
		worker:        make([]*worker.ChatWorker, workerPoolNum),
		handler:       handler,
		wg:            &wg,
	}
}

func (r *ChatProcessor) Start(ctx context.Context) {
	// 내부적으로 worker goroutine N개 실행
	// worker는 채널에서 메시지를 받아 처리
	// ctx.Done()이 닫히면 graceful shutdown

	for i := 1; i <= r.workerPoolNum; i++ {
		// new worker, 채널 공통 주입, wg 주입
		r.wg.Add(1)
		w := worker.NewChatWorker(ctx, r.wg, r.msgChan, r.handler, i)
		r.worker = append(r.worker, w)
		w.Recv()
	}

}

func (r *ChatProcessor) Submit(msg entity.Message) error {
	// 	메시지를 내부 채널에 전달
	// context가 cancel된 경우 에러 반환
	// 채널이 가득 찬 경우 블로킹 or 에러 처리 선택 구현
	for {
		select {
		case r.msgChan <- msg: // 채널에 데이터를 넣을 수 있다면
		default:
			return errors.New("submit error")
		}
	}
}

func (r *ChatProcessor) Stop() {
	// 모든 worker가 정상 종료되도록 처리
	// goroutine leak 없어야 함

	go func() {
		// 모든 worker가 종료 될때 까지 대기
		r.wg.Wait()
		// 모든 worker가 종료되었다면 채널 닫기
		close(r.msgChan)
	}()

	for idx := range r.worker {
		w := r.worker[idx]
		w.Stop()
	}

}
