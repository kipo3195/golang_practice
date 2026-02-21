package usecase

import (
	"fmt"
	"message/repository"
)

// 2. 서비스 구조체
type MessageService struct {
	Repo repository.MessageRepository
}

// 3. 테스트 대상 함수 (SendMessage)
func (s *MessageService) SendMessage(userID int, content string) error {
	// TODO: 비즈니스 로직 구현

	if content == "" {
		// - content가 ""이면 에러 반환
		return fmt.Errorf("content is empty.")
	}

	if len(content) > 100 {
		// - content가 100자 초과하면 에러 반환
		return fmt.Errorf("content length is exceed.")
	}

	// - 통과하면 s.Repo.Save() 호출
	return s.Repo.Save(userID, content)
}
