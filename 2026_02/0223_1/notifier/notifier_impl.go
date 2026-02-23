package notifier

type notifierImpl struct {
}

func NewNotifierImpl() Notifier {
	return &notifierImpl{}
}

func (r *notifierImpl) Send(message string) error {

	// 메시지를 처리하거나, 종료되거나

	//.. 처리 로직
	return nil
}
