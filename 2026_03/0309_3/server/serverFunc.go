package server

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

func ServerA(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, in chan<- string) {
	defer wg.Done()
	min := 0
	max := 2000

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond

	select {
	case <-ctx.Done():
		return
	case <-time.After(sleepDuration):
		select {
		case <-ctx.Done():
			return
		case in <- "serverA":
		}
	}
}

func ServerB(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, in chan<- string) {
	defer wg.Done()
	min := 0
	max := 2000

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond

	select {
	case <-ctx.Done():
		return
	case <-time.After(sleepDuration):
		select {
		case <-ctx.Done():
			return
		case in <- "serverB":
		}
	}
}

func ServerC(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, in chan<- string) {
	defer wg.Done()
	min := 0
	max := 2000

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(max-min+1) + min
	sleepDuration := time.Duration(n) * time.Millisecond

	select {
	case <-ctx.Done():
		return
	case <-time.After(sleepDuration):
		select {
		case <-ctx.Done():
			return
		case in <- "serverC":
		}
	}
}
