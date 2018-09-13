package main

import (
	"testing"
	"sync"
	"fmt"
)

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i:=0; i<10; i++{
		go DoWorkWithReplicatedRequest(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()
	fmt.Printf("receive an answer from #%v\n", firstReturned)
}