package main

import (
	"log"
	"sync"
	"test/counter"
)

func main() {

	c := counter.NewCounter(0)

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 1; i <= 1000; i++ {
		go func() {
			defer wg.Done()
			c.AddCounter()
		}()
	}

	wg.Wait()
	log.Println("result :", c.GetCounter())
	log.Println("작업 종료 ")

}
