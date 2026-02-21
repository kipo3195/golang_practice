package mocking

import "fmt"

type UserUsecase struct {
	repository UserRepository
}

// 인사 메시지를 반환하는 함수
func (s *UserUsecase) GetWelcomMessage(userId int) string {
	// 실제 DB에 연결하지 않고 추상화된 UserRepository에 연결된 척 테스트가 필요한 상황.
	name, err := s.repository.GetUserName(userId)
	if err != nil {
		return fmt.Sprintf("guest")
	}
	return fmt.Sprintf("Welcome, %s!", name)
}
