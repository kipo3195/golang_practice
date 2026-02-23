package main

import (
	"test/dispatcher"
	"test/notifier"
)

// https://gemini.google.com/share/3a187ae76fb5

// 1. 요구 사항
// 인터페이스 정의: Notifier 인터페이스를 만드세요. 이 인터페이스는 Send(message string) error 메서드를 가집니다.
// Dispatcher 구조체: 여러 개의 Notifier를 등록하고 관리하는 Dispatcher를 구현하세요.
// 동시성 처리: Dispatcher의 Broadcast(message string) 메서드는 등록된 모든 Notifier에게 고루틴을 사용하여 메시지를 동시에 전달해야 합니다.
// 제한 시간(Timeout): 각 알림 발송은 최대 2초 이내에 완료되어야 하며, 전체 Broadcast 작업은 모든 고루틴이 종료될 때까지 기다려야 합니다.
// 에러 핸들링: 개별 알림 발송 중 에러가 발생하더라도 다른 알림 발송에 영향을 주지 않아야 합니다.

// 2. 제약 조건
// sync.WaitGroup 또는 channel을 사용하여 모든 고루틴의 종료를 제어하세요.
// context 패키지를 사용하여 타임아웃을 처리하세요.

func main() {

	dispatcher := dispatcher.NewDispatcher()

	for i := 1; i <= 5; i++ {
		dispatcher.RegistNotifier(notifier.NewNotifierImpl())
	}

	dispatcher.BroadCast("test")

}
