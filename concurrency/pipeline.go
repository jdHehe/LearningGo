package main

import (
	"sync"
	"fmt"
)

func gen(nums ... int) <- chan int{
	out := make(chan int)
	go func() {
		for _, n := range nums{
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(done chan int, in <- chan int) <- chan int {
	out := make(chan int)
	go func() {
		// 'for range' can get data out from chanel
		defer close(out)
		for n := range in{
			select {
			case out <- n*n:
			case <-done:
				return
			}
		}
		// remember to close useless channel
	}()
	return  out
}

func merge(done <-chan int, cs ... <-chan int) <-chan int{
	//  <- chan 只能接受值， chan <- 只能发送值

	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <- chan int) {
		for n := range c{
			select {
			case out <- n:
			case <-done:
			}
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _,c := range cs{
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main(){
	//
	in := gen(2,3,4,5,6,7,8)

	done := make(chan int)
	c1 := sq(done, in)
	c2 := sq(done, in)

	defer close(done)
	for n := range merge(done, c1, c2){
		fmt.Println(n)
	}


	//out := merge(done, c1, c2)
	//fmt.Println(<-out)
}
