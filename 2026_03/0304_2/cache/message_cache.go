package cache

import (
	"context"
	"fmt"
	"sync"
	"test/entity"
	"time"
)

type MessageCache struct {
	messageMap map[string]entity.Message

	// 문제점: sync.Mutex나 sync.RWMutex는 절대로 복사되면 안 됩니다.
	mu     sync.RWMutex
	cancel context.CancelFunc
}

func NewMessageCache(cancel context.CancelFunc) *MessageCache {
	return &MessageCache{
		messageMap: make(map[string]entity.Message),
		cancel:     cancel,
	}
}

func (r *MessageCache) Set(key string, msg entity.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messageMap[key] = msg
}

func (r *MessageCache) Get(key string) (entity.Message, bool) {
	// 뺄때 체크
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, exists := r.messageMap[key]
	if !exists {
		return entity.Message{}, false
	}

	// 시간이 지났으면 빈 구조체, 아니라면 value
	exp := value.Exp
	now := time.Now()
	if now.After(exp) {
		return entity.Message{}, false
	}

	return value, true
}

func (r *MessageCache) StartCleanUp(ctx context.Context, s time.Duration) {

	// 클로저, StartCleanUp함수가 종료되는 즉시 ticker 멈춤.. 그러므로 go루틴 안에서 생성, 해제되어야함
	// "고루틴 안에서 사용할 자원(Ticker, Timer 등)의 수명 주기가 고루틴과 같다면, 반드시 고루틴 안에서 생성하고 해제(defer)해야 합니다."
	// ticker := time.NewTicker(s)
	// defer ticker.Stop()

	// 얘가 별도 고루틴으로 돈다면
	go func() {

		ticker := time.NewTicker(s)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.mu.Lock()
				fmt.Println("[cleanup]")
				// 여기서 만료된거 삭제
				// lock 획득 필요함

				// 삭제 대상 키 수집 (맵 순회 중 삭제는 안전함)
				// 로깅을 위해서라면 별도 슬라이스 생성 후 한번만 로깅하고 다시 for문을 통해 delete함수 호출
				// lock 내부에서는 출력은 되도록이면 하지않을것..
				for k := range r.messageMap {
					result := r.messageMap[k]
					exp := result.Exp

					now := time.Now()
					if now.After(exp) {
						// 삭제
						fmt.Println("[cleanup] key :", k)
						delete(r.messageMap, k)
					}
				}
				r.mu.Unlock()
			case <-ctx.Done():
				// 즉시 종료
				return
			}
		}
	}()
}

func (r *MessageCache) StopCleanUp() {
	r.cancel()
}
