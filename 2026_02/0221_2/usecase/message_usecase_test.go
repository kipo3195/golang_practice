package usecase

import (
	"fmt"
	"strings"
	"testing"
)

// --- 테스트 코드 ---

// 4. Mock 객체 정의
type MockMessageRepo struct {
	// TODO: 필요한 필드가 있다면 정의 (예: 호출 여부 확인용)
}

// repository의 로직은 어떠한가
func (m *MockMessageRepo) Save(userID int, content string) error {

	if userID < 0 {
		return fmt.Errorf("userId is not negative.")
	}

	return nil
}

// testSendMessage를 호출
func TestSendMessage(t *testing.T) {
	// 5. Table-Driven Test 구성
	tests := []struct {
		name    string
		userID  int
		content string
		wantErr bool
	}{
		// TODO: 테스트 케이스 추가
		// - 성공 케이스
		{"성공 케이스", 1, "hello", false},
		// - 내용 빈 값 케이스
		{"내용 빈 값 케이스", 1, "", true},
		// - 100자 초과 케이스
		{"100자 초과 케이스", 1, strings.Repeat("a", 100), true},
		// - 음수 에러 케이스
		{"음수 에러 케이스", -1, "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: 서비스 객체 생성 및 Mock 주입 (Dependency Injection)
			mock := &MockMessageRepo{}
			service := MessageService{Repo: mock}

			// TODO: SendMessage 호출 및 결과 검증
			result := service.SendMessage(tt.userID, tt.content)

			// err가 return 되어야할 케이스인데 result 가 nil인 경우
			if tt.wantErr && result == nil {
				t.Errorf("테스트 코드 작성시 예상했던 바와 다름.")
			} else if !tt.wantErr && result != nil {
				t.Errorf("테스트 코드 작성시 예상했던 바와 다름.")
			}

		})
	}
}
