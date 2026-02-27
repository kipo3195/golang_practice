package worker

import (
	"chat/entity"
	"chat/handler"
	"context"
	"sync"
)

type ChatWorker struct {
	ctx     context.Context
	wg      *sync.WaitGroup
	msgChan chan entity.Message
	handler handler.Handler
	idx     int
}

// 공통의 채널을 주입하여 worker 생성
func NewChatWorker(ctx context.Context, wg *sync.WaitGroup, msgChan chan entity.Message, handler handler.Handler, idx int) *ChatWorker {
	return &ChatWorker{
		ctx:     ctx,
		wg:      wg,
		msgChan: msgChan,
		handler: handler,
		idx:     idx,
	}
}

func (r *ChatWorker) Recv() {

	for {
		select {
		case <-r.ctx.Done():
			r.wg.Done()
			return
		case msg := <-r.msgChan:
			r.handler(r.ctx, msg, r.idx)
		}
	}

}

// Stop 함수
// wg.Done() 호출.
func (r *ChatWorker) Stop() {
	r.wg.Done()
}
