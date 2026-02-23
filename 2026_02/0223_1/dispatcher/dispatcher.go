package dispatcher

import (
	"context"
	"test/notifier"
)

type Dispatcher struct {
	messageChan chan string
	Notifier    []notifier.Notifier
	ctx         context.Context
}

func NewDispatcher(ctx context.Context, messageChan chan string) *Dispatcher {
	return &Dispatcher{
		messageChan: messageChan,
		ctx:         ctx,
	}
}

func (r *Dispatcher) RegistNotifier(notifier notifier.Notifier) {
	r.Notifier = append(r.Notifier, notifier)
}

func (r *Dispatcher) Broadcast(message string) {

	defer close(r.messageChan)

	for {
		select {
		case msg := <-r.messageChan:
			for _, n := range r.Notifier {
				go n.Noti(msg)
			}
		case <-r.ctx.Done():
			return
		}
	}

}
