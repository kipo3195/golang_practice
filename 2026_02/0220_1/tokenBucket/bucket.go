package tokenbucket

import (
	"context"
	"time"
)

type TokenBucket struct {
	TokenChan chan struct{}
	cancel    context.CancelFunc
}

func NewTokenBucket(cap int, interval time.Duration) TokenBucket {

	ctx, cancel := context.WithCancel(context.Background())

	tokenChan := make(chan struct{}, cap)

	t := TokenBucket{
		TokenChan: tokenChan,
		cancel:    cancel,
	}

	go t.refill(interval, ctx)

	return t
}

func (r *TokenBucket) refill(interval time.Duration, ctx context.Context) {

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			// 기존 로직
			// r.TokenChan <- struct{}{}
			// 특정 주기마다 채널에 이벤트 전송
			// [주의] 채널은 꽉차면 더이상 데이터를 넣지못하고 소비해 줄때까지 멈춤.
			// ctx cancel이 호출 되더라도..
			// 그러므로 아래 처럼 현재 들어갈 수 있나? 를 체크해야함.
			select {
			case r.TokenChan <- struct{}{}:
				// 토큰 충전 성공
				// r.TokenChan <- struct{}{} -> 이렇게 하면 중복임...위에 상태 체크 + 넣는 행위 같이 하나의 원자적(Atomic)인 동작으로 묶임.
			default:
				// 바구니가 이미 꽉 참 (넘침 방지)
				continue
			}

		case <-ctx.Done():
			// 외부에서 종료 호출
			return
		}
	}

}

func (r *TokenBucket) Allow() bool {

	select {
	case _, ok := <-r.TokenChan: // 꺼낼 수 있는 상태 + 꺼냄
		// 채널이 닫힌 상태일때는 nil이므로 false
		if !ok {
			return false
		} else {
			return true
		}
	default: // 꺼낼 수 없는 상태
		return false
	}
}
