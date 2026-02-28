package worker

import (
	"chat/entity"
	"chat/handler"
	"context"
	"log"
	"sync"
)

type ChatWorker struct {
	ctx     context.Context
	msgChan chan entity.Message
	handler handler.Handler
	idx     int
	Wg      *sync.WaitGroup
}

// 공통의 채널을 주입하여 worker 생성
func NewChatWorker(ctx context.Context, wg *sync.WaitGroup, msgChan chan entity.Message, handler handler.Handler, idx int) *ChatWorker {
	return &ChatWorker{
		ctx:     ctx,
		msgChan: msgChan,
		handler: handler,
		idx:     idx,
		Wg:      wg,
	}
}

func (r *ChatWorker) Recv() {

	for {
		select {
		case <-r.ctx.Done():
			log.Printf("%d worker 종료\n", r.idx)
			r.Wg.Done()
			return
		case msg, ok := <-r.msgChan:
			if !ok {
				r.Wg.Done()
				return
			}
			r.handler(r.ctx, msg, r.idx)
		}
	}

}
