package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// https://gemini.google.com/share/4dcf34e47acd
// 문제 1: 제한된 시간 내 작업 완료 (Timeout & Context)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	longRunningTask(ctx)

}

func longRunningTask(ctx context.Context) {

	// 랜덤한 시간
	min := 1000
	max := 5000

	// 0이상 ()값 미만의 값
	n := rand.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond

	fmt.Println("sleepDuration :", sleepDuration)

	select {
	case <-ctx.Done():
		fmt.Println("Task timed out")
	case <-time.After(sleepDuration):
		fmt.Println("Task completed successfully")

	}
}

// 별도 고루틴에서 무거운 작업 수행 -> 작업이 끝나면 선언한 채널에 종료신호 (닫아버리기)
// func longRunningTask(ctx context.Context) {
//     done := make(chan struct{}) // 작업 완료 신호를 보낼 채널

//     go func() {
//         // 실제 무거운 작업 수행
//         n := rand.Intn(4001) + 1000
//         time.Sleep(time.Duration(n) * time.Millisecond)
//         close(done)
//     }()

//     select {
//     case <-ctx.Done():
//         fmt.Println("Task timed out:", ctx.Err()) // context.DeadlineExceeded 출력 가능
//     case <-done:
//         fmt.Println("Task completed successfully")
//     }
// }
