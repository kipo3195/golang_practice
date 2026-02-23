package dispatcher

import (
	"context"
	"sync"
	"test/notifier"
	"time"
)

type Dispatcher struct {
	notifiers []notifier.Notifier
}

func NewDispatcher() Dispatcher {
	return Dispatcher{
		notifiers: make([]notifier.Notifier, 0),
	}
}

func (r *Dispatcher) RegistNotifier(n notifier.Notifier) {
	r.notifiers = append(r.notifiers, n)
}

func (r *Dispatcher) BroadCast(message string) {

	var wg sync.WaitGroup

	for _, n := range r.notifiers {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		wg.Add(1)
		go func() {
			n.Send(message)
		}()
	}

	wg.Wait()

}
