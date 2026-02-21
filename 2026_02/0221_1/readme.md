1. go에서 주로 사용되는 테스트에 대해 예제 코드를 작성한다.
2. 테스트 유형에는 unitTest, mocking, benchmark가 있다.
3. unitTest 디렉터리에서는 함수가 의도한 바 대로 동작하는지 테스트 하는 예제 코드를 작성한다.
4. mocking 디렉터리에서는 DB나 외부 API를 호출할때 실제 연결 대신 가짜 객체를 생성하여 interface를 활용한 의존성 주입으로 테스트하는 예제 코드를 작성한다.
5. benchmark 디렉터리에서는 testing.B를 통해 특정 로직의 실행속도나 메모리 할당량을 측정하는 예제 코드를 작성한다. 

unitTest는 Table-Driven Test 즉, 정상적인 상황과 예외 상황을 충분히 고려하여 작성하여 테스트 하는 것이 중요하다.
mocking에서는 테스트 하고자 하는 비즈니스 로직은 DB, API의 interface에 의존해야 한다.
benchmark는 ns/op, B/op, allocs/op 측정에 따른 분석으로 성능에 대한 효율을 최대한으로 끌어올리는 것이 중요하다
