package chat

import (
	"chat/entity"
	"chat/handler"
	"chat/processor"
	"chat/worker"
	"context"
	"errors"
	"log"
	"sync"
)

type ChatProcessor struct {
	workerPoolNum int
	msgChan       chan entity.Message
	worker        []*worker.ChatWorker
	cancel        context.CancelFunc
	handler       handler.Handler
	wg            *sync.WaitGroup
}

func NewChatProcessor(workerPoolNum int, handler handler.Handler, cancel context.CancelFunc) processor.Processor {
	var wg sync.WaitGroup
	return &ChatProcessor{
		workerPoolNum: workerPoolNum,
		msgChan:       make(chan entity.Message, 100),
		handler:       handler,
		cancel:        cancel,
		wg:            &wg,
	}
}

func (r *ChatProcessor) Start(ctx context.Context) {

	// worker init
	workers := make([]*worker.ChatWorker, 0, r.workerPoolNum)
	for i := 1; i <= r.workerPoolNum; i++ {
		r.wg.Add(1)
		w := worker.NewChatWorker(ctx, r.wg, r.msgChan, r.handler, i)
		workers = append(r.worker, w)
		go w.Recv()
	}

	r.worker = workers
}

func (r *ChatProcessor) Submit(msg entity.Message) error {
	// 	메시지를 내부 채널에 전달
	// context가 cancel된 경우 에러 반환
	// 채널이 가득 찬 경우 블로킹 or 에러 처리 선택 구현

	// 해결 안 된 것	설명
	// processor가 Stop 중인지	❌
	// channel이 close 되었는지	❌
	// ctx cancel 상태인지	❌
	// graceful shutdown 보장	❌

	select {
	case r.msgChan <- msg: // 채널에 데이터를 넣을 수 있다면
	//case r.ctx.Done() : 로직 필요함.. 그래야 Stop에서 r.cancel()이 호출 됬을때 더이상 채널로 보내지 않도록 처리 할 수 있음
	default:
		log.Println("채널에 데이터를 넣을 수 없음")
		return errors.New("submit error")
	}
	return nil
}

func (r *ChatProcessor) Stop() {
	// 모든 worker가 정상 종료되도록 처리
	// goroutine leak 없어야 함

	// 순서가 중요 하다는데
	r.cancel()

	// 채널을 먼저 닫으면 Recv 로직은 상관없지만 (ok로 체크)
	// send하는쪽에서는 메시지를 보내는 순간에 닫혀있는 상태가 되므로.. ctx로 send와 recv에게 모두 취소신호를 준 뒤 send 하지 못하게 + recv하지 못하게 한 다음
	//  채널을 닫는것이 바람직
	close(r.msgChan)
	r.wg.Wait()

	log.Println("모든 고루틴 종료 ")
}
