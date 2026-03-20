package produce

import (
	"context"
	"fmt"
	"math/rand"
	"test/entity"
)

func Producer(ctx context.Context, imgNum int) chan entity.Job {

	imageChan := make(chan entity.Job)

	go func(imgNum int) {
		// 쓰는 쪽에서 닫아야하므로
		defer close(imageChan)
		for i := 1; i <= imgNum; i++ {
			r := rand.Intn(5) + 1
			select {
			case <-ctx.Done():
				return
			case imageChan <- entity.Job{
				ID:         i,
				TaskName:   fmt.Sprintf("%d - ImageTask", i),
				Complexity: r,
			}:
			}
		}

	}(imgNum)

	return imageChan
}
