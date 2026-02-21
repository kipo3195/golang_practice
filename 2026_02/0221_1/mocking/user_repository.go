package mocking

type UserRepository interface {
	GetUserName(userId int) (string, error)
}
