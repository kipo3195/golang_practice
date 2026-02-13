package main

import (
	"log"
	"sync"
	"time"
)

// N개의 작업을 M개의 worker(비동기 고루틴)가 처리하는 예제
func main() {

	// wg init
	var wg sync.WaitGroup

	// n개의 작업과 m개의 worker 정의
	job := 10
	worker := 3

	// 작업 전달을 위한 채널 생성
	jobChan := make(chan int)

	log.Println("작업 시작 ")
	// m개의 workerpool 생성
	for i := 1; i <= worker; i++ {
		wg.Add(1)
		go print(&wg, jobChan, i)
	}

	// n개의 작업 실행
	for i := 1; i <= job; i++ {
		jobChan <- i
	}

	close(jobChan)

	wg.Wait()
	log.Println("작업 종료 ")

}

// 비동기 로직 처리함수
func print(wg *sync.WaitGroup, jobChan <-chan int, i int) {

	defer wg.Done()

	for value := range jobChan {
		log.Printf("%d 번째 worker 수행. 출력값 : %d", i, value)
		time.Sleep(1 * time.Second)
	}

}
