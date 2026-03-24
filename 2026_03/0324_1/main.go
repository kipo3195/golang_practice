package main

import (
	"container/heap"
	"fmt"
)

type Process struct {
	pid      int
	priority int
}

// heap을 사용해서 처리하기 위한 interface구현
type PriortyQueue []*Process

func (pq PriortyQueue) Len() int {
	return len(pq)
}

// 짧고 작은 <
func (pq PriortyQueue) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return pq[i].pid < pq[j].pid
	}
	return pq[i].priority < pq[j].priority
}

func (pq *PriortyQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq *PriortyQueue) Push(x any) {
	item := x.(*Process)
	*pq = append(*pq, item)
}

func (pq PriortyQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func main() {
	pq := &PriortyQueue{}
	heap.Init(pq)

	heap.Push(pq, &Process{pid: 1, priority: 10})
	heap.Push(pq, &Process{pid: 2, priority: 5})
	heap.Push(pq, &Process{pid: 3, priority: 15})
	heap.Push(pq, &Process{pid: 4, priority: 5})

	fmt.Println("--- 처리순서 ---")

	for pq.Len() > 0 {
		p := heap.Pop(pq).(*Process)
		fmt.Printf("[priority %d] %d\n", p.priority, p.pid)
	}
	// 예상 p2, p4, p1, p3

}
