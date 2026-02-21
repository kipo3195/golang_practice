package mocking

import (
	"fmt"
	"testing"
)

// DB 연결대신 사용될 repository - usecase에서 연결될 repository가 추상화되있는 것이 중요함.
// usecase는 의존성 주입을 통해 추상화된 user_repository를 알고있음. user_repository는 인터페이스로서,
// 내부에 GetUserName이라는 함수를 구현해야함.
// 여기서 생성할 임의의 repository (MockUserRepository)는 당연히 GetUserName() 함수를 구현해야함.
type MockUserRepository struct {
	FackName   string
	ShouldFail bool
}

// UserRepository interface의 함수 구현
func (m *MockUserRepository) GetUserName(userId int) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("DB 에러 발생")
	}
	return m.FackName, nil
}

// 아마 usecase의 GetWelcomMessage함수를 테스트
func TestGetWelcomMessage(t *testing.T) {
	t.Run("성공시 사용자 이름 포함", func(t *testing.T) {
		// DB 연결 repo가 아닌 테스트용 repo. MockUserRepository시 ShouldFail을 true로 주지 않으면 사용자 이름을 반환.
		mock := &MockUserRepository{FackName: "Gopher"}
		// 실제 usecase
		usecase := &UserUsecase{repository: mock}

		result := usecase.GetWelcomMessage(1)
		expected := "Welcome, Gopher!"

		if result != expected {
			t.Errorf("got %s, want %s", result, expected)
		}
	})

	t.Run("실패시 DB 에러 발생", func(t *testing.T) {
		mock := &MockUserRepository{ShouldFail: true}
		usecase := &UserUsecase{repository: mock}

		result := usecase.GetWelcomMessage(1)
		expected := "guest"

		if result != expected {
			t.Errorf("got %s, want %s", result, expected)
		}
	})

}
