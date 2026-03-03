package worker

import (
	"context"
	"math/rand"
	"sync"
	"test/dto"
	"time"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, jobName string, resultChan chan<- dto.Result) {

	select {
	case <-ctx.Done():
		wg.Done()
		return
	default:
		min := 100
		max := 800
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := r.Intn(max-min+1) + min
		time.Sleep(time.Duration(n) * time.Millisecond)
		var err string
		if jobName == "fail" {
			err = "error"
		} else {
			err = ""
		}
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case resultChan <- dto.Result{
			Value: jobName,
			Err:   err,
		}:
		}
	}

}
