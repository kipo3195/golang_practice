package main

import (
	"context"
	"fmt"
	"test/cache"
	"test/entity"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := cache.NewMessageCache(cancel)

	// 청소기 시작 1초 간격
	cache.StartCleanUp(ctx, 1*time.Second)

	cache.Set("A", entity.Message{
		Value: "A",
		Exp:   time.Now().Add(2 * time.Second),
	})
	cache.Set("B", entity.Message{
		Value: "B",
		Exp:   time.Now().Add(5 * time.Second),
	})

	// A 출력 확인
	value, exist := cache.Get("A")
	fmt.Println("Get A:", value, exist)

	time.Sleep(3 * time.Second)

	value, exist = cache.Get("A")
	fmt.Println("Get A (after 3 second):", value, exist)
	value, exist = cache.Get("B")
	fmt.Println("Get B (after 3 second):", value, exist)

	cache.StopCleanUp()
	fmt.Println("작업 종료.")
}
