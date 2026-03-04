package processor

import (
	"context"
	"sync"
	"time"
)

func Processor(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// in의 모든 데이터를 소진하고나면 for루프 종료

			// for v := range in 는 in에 데이터가 들어와야 로직이 수행됨.
			// 어떤 이유로 인해 in에는 채널이 들어오지 않는 상태
			// 그럼 in로직에 걸려있게 됨 -> 아래로 내려가지 않으므로 ctx.Done()을 탈 수 없음.
			// for v := range in {
			// 	select {
			// 	case <-ctx.Done():
			// 		// 데이터를 꺼내 쓰려고하는 시점에 타임아웃 혹은 취소 이벤트를 수신했는가 점검
			// 		return
			// 	default:

			// 		time.Sleep(10 * time.Millisecond)
			// 		sqaure := v * v
			// 		// 데이터를 채널에 전송하기 전에 ctx.Done()에 의해 타임아웃, 취소 이벤트를 수신했는가 점검
			// 		select {
			// 		case <-ctx.Done():
			// 			return
			// 		case out <- sqaure: // 보낼 수 있는 상태인가 + 보내기 동시 처리
			// 		}
			// 	}
			// }

			// 개선 ====>
			// 루프를 반복적으로 돌면서 취소를 즉각 체크 (채널에 데이터가 들어오는지도 지속적으로 체크 )
			for {
				select {
				case <-ctx.Done():
					// 취소 신호에 대한 즉각적인 반응이 가능함
					return
				case v, ok := <-in:
					if !ok {
						// 채널이 닫힌 상태이므로 더이상 처리하지 않음
						return
					}

					// 비즈니스 로직 수행
					time.Sleep(10 * time.Millisecond)
					sqaure := v * v
					select {
					case <-ctx.Done():
						// 데이터를 보낼 수 없는 상태 = 타임아웃, 취소
						return
					case out <- sqaure:
						// 데이터를 보낼 수 있는 상태 + 보냄
					}
				}
			}

		}()
	}

	// 모든 wg.Done()이 호출되고나면 채널 닫기
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
