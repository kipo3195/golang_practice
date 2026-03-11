package entity

type Task struct {
	DeviceID int // 이 ID를 기준으로 워커가 결정됨
	Payload  string
}
