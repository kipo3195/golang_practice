package notifier

import (
	"errors"
	"sync/atomic"
	"time"
)

type MockNotifier struct {
	CallCount  int32
	ShouldFail bool
	Delay      time.Duration
}

// 잘 호출 되는가
// 의도적으로 지연했을때 에러를 잘 뱉는가
func (r *MockNotifier) Send(msg string) error {

	// data race condition방지를 위해 atomic 패키지 활용
	atomic.AddInt32(&r.CallCount, 1)

	if r.Delay > 0 {
		time.Sleep(r.Delay)
	}
	if r.ShouldFail {
		return errors.New("send failed")
	}

	return nil
}
