package sum

import (
	"test/common"
)

// 정수 슬라이스를 받아 각 요소에 +1 을 수행하는 함수를 반환하는 '작업 생성기' 함수
// 고차 함수
// SumWork는 굳이 매개변수를 받을 필요 없음.
func SumWorkOperator() common.Operator {

	// 클로저 형식
	// SumWork는 함수를 리턴하고, 그 함수의 내부 로직은 int를 리턴한다.
	return func(n int) int {
		return n + 1
	}
}
