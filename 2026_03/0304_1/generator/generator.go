package generator

import (
	"context"
)

func Generator(ctx context.Context, s []int) <-chan int {

	// 비동기로 채널에 데이터를 삽입
	in := make(chan int)

	go func() {
		// 고루틴이 끝날때 채널을 닫는게 맞음
		// Generator 함수가 끝날때 닫는다면 고루틴이 데이터를 보내기도 전에 가게 문을 닫아버리는 꼴이 됩니다.
		// -> 외부에서 in 채널을 바라보고있는 루프가 즉시종료됨 + 보내는 쪽은 panic

		defer close(in)
		// 채널은 데이터를 보내는 로직이 닫는게 맞음.

		for _, v := range s {
			// ctx는 채널에 담기 전에 매번 체크하도록 함.
			select {
			case in <- v:
				// 넣을수 있으면 넣고
			case <-ctx.Done():
				// 취소 신호를 받았으니 종료.
				return
			}
		}
	}()

	return in
}
