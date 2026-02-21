package repository

// 1. Repository 인터페이스 (Mocking 대상)
type MessageRepository interface {
	Save(userID int, content string) error
}
