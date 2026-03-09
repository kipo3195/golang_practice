package main

import (
	"context"
	"fmt"
	"sync"
	"test/server"
)

// https://gemini.google.com/share/4dcf34e47acd
// 문제 3: 가장 먼저 도착한 결과 사용 (First Response / Fan-in)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan string)

	var wg sync.WaitGroup

	wg.Add(3)
	go server.ServerA(ctx, cancel, &wg, in)
	go server.ServerB(ctx, cancel, &wg, in)
	go server.ServerC(ctx, cancel, &wg, in)

	select {
	case result := <-in:
		fmt.Println("result :", result)
		cancel()
		close(in)
	}

	wg.Wait()
	fmt.Println("작업 종료")

}
