package main

import (
	"context"

	"./inMemory"
)

// 요구사항
//
//	기능 : `Set(key, value, ttl)`, `Get(key)`, `Delete(key)` 메서드를 구현하세요.
//	동시성 : 여러 고루틴이 동시에 접근해도 데이터 깨짐이 없어야 합니다. (`sync.RWMutex` 활용 추천)
//	만료 처리 (TTL)
//	  - 데이터 저장 시 만료 시간(TTL)을 설정합니다.
//	  - 만료된 데이터는 `Get` 시점에 조회되지 않아야 합니다.
//	핵심 별도의 고루틴(Cleanup Goroutine)이 주기적으로 메모리를 스캔하여 만료된 데이터를 삭제해야 합니다.
//	범용성 (제네릭) : Go 1.18+의 제네릭을 사용하여 다양한 타입을 담을 수 있게 만드세요.
//	인터페이스 : 향후 Redis 등으로 교체 가능하도록 `Cache` 인터페이스를 먼저 정의하세요.

// 20260219 범용성 반영 하지 않은
func main() {
	cache := inMemory.NewCacheInterfaceImpl()
	ctx, cancel := context.WithCancel(context.Background())
	go cache.CleanUp(ctx, 5)
	cancel()
}
