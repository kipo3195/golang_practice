package collector

import "context"

func Collector(ctx context.Context, out <-chan int) []int {

	result := make([]int, 0)

	// processor와 같은이유로 주석처리함 (즉각처리 안됨 )
	// for v := range out {
	// 	select {
	// 	case <-ctx.Done():
	// 		return result
	// 	default:
	// 		result = append(result, v)
	// 	}
	// }

	for {
		select {
		case <-ctx.Done():
			return result
		case v, ok := <-out:
			if !ok {
				return result
			}
			result = append(result, v)
		}
	}

}
