package notifier

import (
	"context"
	"log"
	"sync"
)

type NotifierImpl struct {
	messageChan chan string
	wg          *sync.WaitGroup
	cancel      context.CancelFunc
}

func NewNotifier(wg *sync.WaitGroup, cancel context.CancelFunc) Notifier {

	wg.Add(1)

	return &NotifierImpl{
		wg:          wg,
		messageChan: make(chan string),
		cancel:      cancel,
	}
}

func (r *NotifierImpl) Send(message string) error {
	defer r.wg.Done()
	for {
		select {
		case msg := <-r.messageChan:
			log.Println("msg 처리 : ", msg)
			// 2초 이상 걸리면
			if message == "stop" {
				r.cancel()
				return nil
			}
		}
	}

}

func (r *NotifierImpl) Noti(message string) {
	r.messageChan <- message
}
