package main

import (
	"context"
	"test/producer"
	"time"
)

// 팬아웃 및 배치 로그 플러셔 (Log Flusher)

// 대량의 채팅 메시지가 발생하는 메신저 서버에서, 모든 메시지를 발생할 때마다 DB에 쓰면 부하가 너무 큽니다.
//  1) 메시지를 채널에 모았다가
//  2) 3초마다 한 번씩(Ticker)
//  3) 버퍼가 5개 쌓일 때마다 한꺼번에 DB에 저장하는 '배치 플러셔'를 고루틴으로 구현해야 합니다.

// 메시지 수집: Message 구조체(내용 포함)를 받는 채널을 만듭니다.
// 배치 조건: * 3초 주기로 동작하는 time.Ticker를 사용합니다.

// 메시지가 5개 쌓이면 즉시 저장합니다.

// 동시성 제어: context.WithTimeout을 사용하여 전체 프로세스를 10초 후에 안전하게 종료(Graceful Shutdown) 시키세요.

// 리소스 정리: 종료 시 반드시 ticker.Stop()을 호출하고, 남아있는 메시지가 있다면 모두 출력(저장)하고 종료해야 합니다.

// producer
//    메시지 수집 채널에 데이터 in
// entity
//    Mesage : 인덱스 + 메시지 입니다.
//    FlushEntity : 내부에 값, isFlush, Message 구조체의 배열 -> flush 되고나면 초기화.
// batch
//    batch : 메시지 수집, 3초 주기로 동작

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 메시지 수집, 수집 채널 반환
	msgChan := producer.Produce(ctx, 40)

	// 배치 작업 수행 로직 처리

}
