package chatroom

import "sync"

// 요구사항
// ChatRoom struct를 만들어라.
// 이 채팅방은 최근 100개의 메시지만 유지해야 한다.
// 메시지가 100개를 초과하면 가장 오래된 메시지를 삭제한다.
// 동시성 안전해야 한다. (goroutine 여러개가 동시에 추가 가능)
// ringBuffer로 구현

type ChatRoom struct {
	buffer []string
	mu     sync.RWMutex
	max    int
	start  int
	count  int
}

func NewChatRoom(max int) *ChatRoom {
	if max <= 0 {
		panic("max must be greater than 0")
	}
	return &ChatRoom{
		buffer: make([]string, max), // 하나의 배열만 생성함.
		max:    max,
	}
}

func (c *ChatRoom) AddMessage(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 0~99 까지 0~99인덱스
	// 100 ~ 199 까지 다시 0 ~ 99 인덱스
	idx := (c.start + c.count) % c.max
	c.buffer[idx] = msg

	if c.count < c.max {
		// max의 수까지만 늘어남. 이후 start값 증가 + 모듈러 연산을 통해 인덱스를 계속 순환하게 함
		c.count++
	} else {
		// start를 옮김
		c.start = (c.start + 1) % c.max
	}
}

// 인덱스 변경 예시 max 3
// start 0, count 0 -> idx = 0이므로 c.buffer[0]에 msg 삽입, 이후 c.count++
// start 0, count 1 -> idx = 1이므로 c.buffer[1]에 msg 삽입, 이후 c.count++
// start 0, count 2 -> idx = 2이므로 c.buffer[2]에 msg 삽입, 이후 c.count++
// start 0, count 3 -> idx = 0이므로 c.buffer[0]에 msg 삽입, 이후 start(0) + 1 % 3에 의해 start값 1
// start 1, count 3 -> idx = 1이므로 .... 반복

func (c *ChatRoom) GetRecentMessages() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]string, c.count)
	for i := 0; i < c.count; i++ {
		result[i] = c.buffer[(c.start+i)%c.max]
	}
	return result
}

// 조회 예시
// start 1, count 3일때 조회 순서는 buffer 인덱스 1 (start), 2, 0 순서로 조회되어야함
// count는 이미 max값만큼 이므로 result 배열의 길이를 count만큼
// i는 0 ~ 2까지 순회
// result[0] = c.buffer[1 % 3] [1]
// result[1] = c.buffer[2 % 3] [2]
// result[2] = c.buffer[3 % 3] [0]
