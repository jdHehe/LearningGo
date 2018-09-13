package main

import "fmt"

func Fibnacci(n int) <-chan int{
	result := make(chan int)
	fmt.Printf("Fibnacci（%d)", n)
	go func() {
		//goroutine 是tasks
		defer close(result)
		if n <= 2 {
			result <- 1
			return
		}
		result <- <-Fibnacci(n-1) + <-Fibnacci(n-2)
	}()
	// goroutine之后的称为continuation
	return result
}

func main(){
	fmt.Printf("fib(4) = %d", <-Fibnacci(10))
}