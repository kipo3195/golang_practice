package main

import (
	"fmt"
	"sync"
)

// SendWelcomeMessages는 사용자 목록을 받아 비동기로 인사를 출력합니다.
func SendWelcomeMessages(users []string) {
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		// TODO: 고루틴과 클로저를 활용하여 구현하세요.
		// 주의: 모든 사용자가 자신의 이름을 정확히 출력받아야 합니다.
		go func(u string) {
			defer wg.Done()
			fmt.Printf("Welcome, %s!\n", u)
		}(user)
	}

	wg.Wait()
}

func main() {
	users := []string{"Alice", "Bob", "Charlie", "Dave"}
	SendWelcomeMessages(users)
}
