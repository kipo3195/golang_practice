package chatroom

import "sync"

// 요구사항
// ChatRoom struct를 만들어라.
// 이 채팅방은 최근 100개의 메시지만 유지해야 한다.
// 메시지가 100개를 초과하면 가장 오래된 메시지를 삭제한다.
// 동시성 안전해야 한다. (goroutine 여러개가 동시에 추가 가능)
// slice로 구현

type ChatRoom struct {
	buffer []string
	mu     sync.RWMutex
	max    int
}

func NewChatRoom(max int) ChatRoom {
	return ChatRoom{
		buffer: make([]string, 0, max), // 길이, 용량
		max:    max,
	}
}

func (c *ChatRoom) AddMessage(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.buffer) >= c.max {
		// 1번째 인덱스 부터 끝까지 복사하면되지않을까? -> 전체 복사면 O(N)
		c.buffer = c.buffer[1:]
		// 슬라이스는 backing array를 공유함.
		// 앞 요소는 GC 대상 안됨 (참조가 남아있기 때문)
		// 메모리 오래 유지될 수 있음
		// 계속 밀리면 내부 array 계속 커질 수도 있음 (특정 상황에서)
		// 그리고 논리적으로는 O(n) 복사 발생 가능성 있음.

		// slice는 header 구조체이며
		// 1. pointer (backing array 주소)
		// 2. length
		// 3. capacity
		// 로 구성된다.
		//
		// 대입(b := a)이나 재슬라이싱(a[1:])은
		// backing array는 공유하고 header만 새로 만들어진다.
		// 따라서 한쪽에서 요소를 변경하면 다른 쪽에도 영향을 준다.
		//
		// 완전히 분리하려면 make로 새 backing array를 만들고
		// copy()로 데이터를 복사해야 한다.

	}
	c.buffer = append(c.buffer, msg)
}

func (c *ChatRoom) GetRecentMessages() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 원본 보호
	result := make([]string, len(c.buffer))
	copy(result, c.buffer) // copy하지 않으면 외부에서 슬라이스 수정 가능
	return result
}
