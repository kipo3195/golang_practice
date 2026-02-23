package dispatcher

import (
	"sync/atomic"
	"test/notifier"
	"testing"
	"time"
)

func TestDispatcher_Broadcast(t *testing.T) {

	t.Run("정상적인 케이스 확인", func(t *testing.T) {
		dispatcher := &Dispatcher{}

		m1 := notifier.MockNotifier{}
		m2 := notifier.MockNotifier{}
		dispatcher.RegistNotifier(&m1)
		dispatcher.RegistNotifier(&m2)

		dispatcher.BroadCast("test1")

		if atomic.LoadInt32(&m1.CallCount) != 1 || atomic.LoadInt32(&m2.CallCount) != 1 {
			t.Errorf("Broadcast가 모든 Notifier를 호출하지 않았습니다.")
		}

	})

	t.Run("타임 아웃 케이스 확인", func(t *testing.T) {
		dispatcher := &Dispatcher{}

		m1 := notifier.MockNotifier{Delay: 5 * time.Second}

		dispatcher.RegistNotifier(&m1)
		start := time.Now()
		dispatcher.BroadCast("time out test")
		end := time.Since(start)

		if end > 2500*time.Millisecond {
			t.Errorf("broad cast가 타임아웃 제한을 초과했습니다 %v", end)
		}
	})

}
