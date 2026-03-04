package main

import (
	"context"
	"log"
	"test/collector"
	"test/generator"
	"test/processor"
)

// 당신은 대용량 데이터를 처리하는 시스템을 설계해야 합니다. 데이터 처리 과정은 3단계로 나뉩니다.
// Generator (생성): 숫자를 1부터 N까지 빠르게 생성합니다.
// Processor (가공 - 느림): 생성된 숫자를 받아 복잡한 연산(제곱 + 10ms 지연)을 수행합니다. 이 단계가 병목이므로 여러 개의 고루틴으로 병렬 처리해야 합니다.
// Collector (수집): 처리된 결과를 하나로 모아서 출력합니다.

// Pipeline 구조: 각 단계는 <-chan int (수신 전용 채널)를 반환하는 함수로 구현하세요. 즉, 함수를 호출하면 채널이 튀어나와야 합니다.
// Fan-Out (분산): Processor 단계를 3개의 고루틴이 동시에 수행하도록 만드세요.
// Fan-In (병합): 3개의 Processor가 뱉어내는 채널 3개를 하나의 채널로 합치는 Merge 함수를 작성하세요.

// https://gemini.google.com/share/cbcbc5f82aa9

// close(채널)의 핵심은
// "닫힌 이후에 더 이상 추가 데이터는 못 받지만, 이미 들어온 데이터는 전부 처리된다"

func main() {

	// Generator 숫자를 생성하고 넣은 채널 리턴
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := generator.Generator(ctx, slice)
	out := processor.Processor(ctx, in)
	// result는 슬라이스를 받도록 처리함
	result := collector.Collector(ctx, out)

	for _, v := range result {
		log.Println("value : ", v)
	}

	log.Println("작업 종료")
}
