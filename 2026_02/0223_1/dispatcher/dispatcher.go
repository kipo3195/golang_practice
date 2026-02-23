package dispatcher

import (
	"context"
	"log"
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

type SendError struct {
	NotifierID int   // 어떤 노티파이어인지 식별
	Err        error // 발생한 에러 내용
}

func (r *Dispatcher) BroadCast(message string) []SendError {

	var wg sync.WaitGroup

	// 에러를 수집할 채널 (버퍼 크기를 노티파이어 개수만큼 설정하면 블로킹 방지 가능)
	errChan := make(chan SendError, len(r.notifiers))

	for i, n := range r.notifiers {
		wg.Add(1)
		// defer은 for문의 종료가 아닌 BroadCast 함수 종료시 호출된다.
		// defer wg.Done()
		// defer wg.Done() 는 반드시 실행되는 고루틴 내부에서 호출해야함.

		// 마찬가지로 defer은 for은 종료가 아닌 BroadCast 함수 종료시 호출됨.
		// defer close(resultChan)

		// 각 Notifier를 개별 고루틴으로 실행
		go func(idx int, notifier notifier.Notifier) {
			// 고루틴 종료후 확실히 종료
			defer wg.Done()

			// 알림 발송에 제한시간 2초
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// 결과 수신 채널 생성
			resultChan := make(chan error, 1)

			// 실제 전송 로직 수행
			go func() {
				// 여기서 고루틴이 할당되기도 전에 for루프가 다 돌아버리면 n은 for루프의 가장 마지막 n으로 할당될 가능성이 있음
				// 1.22 버전부터는 개선되었다고하는데..
				// 그래서 for문마다 n을 notifier이라는 로컬 변수에 복사되어 저장됨
				resultChan <- notifier.Send(message)
			}()

			// 결과 대기 및 타임아웃 처리
			select {
			case err := <-resultChan:
				// 로직 수행 완료
				if err != nil {
					log.Println("로직 에러")
					errChan <- SendError{
						NotifierID: idx,
						Err:        err,
					}
				} else {
					log.Println("로직 정상 수행")
				}
			case <-ctx.Done():
				log.Println("타임아웃! ")
				// timeout
			}
		}(i, n)

		// for 루프는 굉장히 빠릅니다. 고루틴이 실제로 실행 준비를 마치고 메모리에 올라가기도 전에 루프는 이미 끝까지 돌아버릴 수 있습니다.
		// 이때 모든 고루틴은 n이라는 동일한 메모리 주소를 바라보고 있습니다.
		// 루프가 끝나면 n에는 슬라이스의 마지막 요소가 담겨 있습니다.
		// 결과적으로, 5명에게 보내려고 했는데 마지막 사람에게만 5번 메시지가 발송되는 대참사가 발생할 수 있습니다. (Go 1.22 버전부터는 이 동작이 개선되었으나, 하위 호환성과 명확한 코딩 습관을 위해 인자로 넘기는 것이 관례입니다.)
		// 이렇게 하면 고루틴이 생성되는 순간의 n 값이 notifier라는 로컬 변수에 복사되어 저장됩니다.
		// 루프가 다음으로 넘어가서 n 값이 바뀌어도,
		// 이미 생성된 고루틴 내부의 notifier 값은 안전하게 보존됩니다.
	}

	// BroadCast는 모든 고루틴이 종료될때까지 대기.
	wg.Wait()

	close(errChan)

	// 채널에 쌓인 에러들을 슬라이스로 변환하여 반환
	var collectedErrors []SendError
	for e := range errChan {
		collectedErrors = append(collectedErrors, e)
	}

	return collectedErrors
}
