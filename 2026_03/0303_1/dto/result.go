package dto

// 결과 구조체 정의: 성공 값(Value)과 에러(Err)를 함께 담을 수 있는 Result 구조체를 정의하고 채널을 통해 이 구조체를 전달하세요.
type Result struct {
	Value string
	Err   string
}
