package collector

import "context"

func Collector(ctx context.Context, out <-chan int) []int {

	result := make([]int, 0)

	for v := range out {
		result = append(result, v)
	}

	return result
}
