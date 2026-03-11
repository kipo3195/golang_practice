package producer

import (
	"fmt"
	"test/entity"
)

func Produce() chan entity.Task {

	// produce의 생성 데이터를 담는 버퍼는 보통 CPU * 2로 처리
	taskChan := make(chan entity.Task, 5)

	go func() {
		// 1~5의 워커가 나누어서 이벤트 처리
		for i := 1; i <= 20; i++ {
			mod := i % 5
			taskChan <- entity.Task{
				DeviceID: mod,
				Payload:  fmt.Sprintf("%d 워커가 처리해야할 payload : %d\n", mod, i),
			}
		}
		close(taskChan)
	}()

	return taskChan
}
