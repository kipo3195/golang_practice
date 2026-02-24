package worker

import (
	"sort"
	"testing"
)

func TestSquareWorker(t *testing.T) {

	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 4, 9, 16, 25}

	worker := SquareWorker()

	results := SquareProcess(input, worker)

	if len(expected) != len(results) {
		t.Errorf("길이가 다릅니다 예상 %d, 결과 %d", len(expected), len(results))
	}

	sort.Ints(results)
	sort.Ints(expected)

	for i := range results {
		if results[i] != expected[i] {
			t.Errorf("인덱스 %d에서 불일치 예상 %d, 결과 :%d", i, expected[i], results[i])
		}
	}

}
