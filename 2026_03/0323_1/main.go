package main

import (
	"container/heap"
	"fmt"
)

type Patient struct {
	name     string
	priority int // 숫자가 높을 수록
}

type PriorityQueue []*Patient

func (pq PriorityQueue) Len() int { return len(pq) }

// 내림차순은 i > j
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Patient)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)      // 길이를 구한다
	item := old[n-1]   // 가장 마지막 항목
	*pq = old[0 : n-1] // 0부터 가장마지막 항목 직전까지 복사
	return item
}

func main() {
	pq := &PriorityQueue{}
	heap.Init(pq)
	// 함수의 인자로 Interface 타입을 받음. Interface는 sort.Interface // Len(), Less(i, j int) bool, Swap(i, j int) 포함,
	// 	Push(x any) // add x as element Len() Pop() any   // remove and return element Len() - 1. 포함
	// 그러므로 Init시점에 매개변수로 던져지는 구조체가 Len(), Less(), Swap(), Len(), Pop()이 구현되어 있어야함.

	// 환자 추가
	heap.Push(pq, &Patient{name: "가벼운 찰과상", priority: 1})
	heap.Push(pq, &Patient{name: "긴급 수술", priority: 10})
	heap.Push(pq, &Patient{name: "독감", priority: 3})

	fmt.Println("--- 치료순서 ---")

	for pq.Len() > 0 {
		p := heap.Pop(pq).(*Patient)
		fmt.Printf("[%d순위] %s\n", p.priority, p.name)
	}

}
