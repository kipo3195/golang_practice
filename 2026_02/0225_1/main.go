package main

import (
	"log"
	"test/common"
	"test/worker"
)

// 클로저 활용: 정수 슬라이스를 받아 각 요소에 +1 을 수행하는 함수를 반환하는 '작업 생성기' 함수를 작성하세요.
// 고루틴 및 채널 활용: 생성된 작업을 여러 개의 고루틴으로 동시에 실행하고, 그 결과물들을 채널을 통해 수집하여 다시 슬라이스로 반환하는 함수를 작성하세요.
// 동기화: 모든 고루틴이 종료될 때까지 기다린 후 채널을 닫아야 하며, 결과 슬라이스의 순서는 상관없습니다.
// 테스트 코드: 작성한 로직이 올바르게 동작하는지 검증하는 단위 테스트(go test) 코드를 포함하세요.

func main() {

	var s []int
	for i := 1; i <= 10; i++ {
		s = append(s, i)
	}

	repo := worker.NewSumWorkerRepository(5)

	op := repo.SumWorkOperator()

	result := common.OperatorProcess(s, op)

	log.Println("결과 : ", result)

}
