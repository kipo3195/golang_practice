// main_test.go
package main

import (
	"log"
	"testing"
)

func TestSendWelcomeMessages(t *testing.T) {
	users := []string{"User1", "User2", "User3"}

	// 이 함수를 실행했을 때 Race Condition이 발생하는지 확인하는 것이 주 목적입니다.
	// 실제 출력값을 캡처하여 검증하려면 조금 더 복잡한 설정이 필요하므로,
	// 여기서는 함수가 데드락 없이 정상 종료되는지와 레이스 감지에 집중하세요.
	SendWelcomeMessages(users)
	log.Println("작업 종료")
}
