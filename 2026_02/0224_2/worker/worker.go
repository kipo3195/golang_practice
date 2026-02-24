package worker

// 어떤 작업이든지 int 값을 던져서 int값을 return할 수 있는 함수
// 타입처럼 사용할 수 있음. (매개변수로 던질 수 도 있고, 반환값으로 사용할 수 있음. )
type Work func(int) int
