package main

import (
	"sync"
	"test/java"
)

const TOTAL_CONNECTIONS = 100

func main() {

	var wg sync.WaitGroup
	java.Connect(&wg)
	wg.Wait()
}
