package notifier

import "log"

type notifierImpl struct {
	messageChan chan string
}

func NewNotifierImpl() Notifier {
	return &notifierImpl{
		messageChan: make(chan string, 1),
	}
}

func (r *notifierImpl) Send(message string) error {

	// 메시지를 처리하거나, 종료되거나
	select {
	case r.messageChan <- message:
		log.Println("msg : ", message)
	}

	return nil
}
