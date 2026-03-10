package producer

import "context"

func Produce(ctx context.Context, dataChan chan<- int, limit int) {

	go func(limit int) {
		// 모든 데이터를 전송 후 채널 닫기
		defer close(dataChan)
		for i := 1; i <= limit; i++ {
			select {
			case dataChan <- i:
			case <-ctx.Done():
				// 취소가 되었는지 매번 점검하는 구조
				return
			}
		}
	}(limit)

}
