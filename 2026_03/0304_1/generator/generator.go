package generator

import (
	"context"
)

func Generator(ctx context.Context, s []int) <-chan int {

	in := make(chan int)

	// 비동기로 채널에 데이터를 삽입
	go func() {
		for _, v := range s {
			in <- v
		}
		close(in)
	}()

	// 데이터를 다 넣은다음 채널은 닫아주기
	// 닫아줘도 사용하는쪽에서는 모두 소진 가능하므로

	return in
}
