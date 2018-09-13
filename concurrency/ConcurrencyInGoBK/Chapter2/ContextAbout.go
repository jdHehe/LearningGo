package main

import (
	"sync"
	"fmt"
	"time"
	"context"
)

// usage of done channel to manage the call-graph
// 1 处产生一个新的context并用WithCancel包装成具有cancellation的context
// 2 处的cancel由main函数调用，当发生print错误的时候进行调用
// 3 处对context对象进一步的封装, 超时的时候它会调用cancel
// 4 处context.Done()调用，并返回结束的原因

// data-bag for context to store or retrieve data
// context.WithValue(context, key, value) : store key-value pair
// context.Value(key) : get the value by key
// 要求： 存储的key-calue数据要求是线程安全的， key的实现要

func main() {

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) // 1
	defer cancel()

	wg.Add(1)
	go func() {
		wg.Done()
		if err := printGreeting_context(ctx); err != nil {
			fmt.Printf(" cannot print greeting:  %v\n", err)
			cancel() // 2
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell_context(ctx); err != nil{
			fmt.Printf(" cannot print farewell: %v \n", err)
			cancel()
		}
	}()

	wg.Wait()
}
func printFarewell_context(ctx context.Context) error {
	farewell, err := genFarewell__context(ctx)
	if err != nil{
		return err
	}
	fmt.Printf(" %s world ! \n", farewell)
	return nil
}
func printGreeting_context(ctx context.Context) error {
	greeting, err := genGreeting_context(ctx)
	if err != nil{
		return err
	}
	fmt.Printf(" %s world ! \n", greeting)
	return nil
}
func genGreeting_context(ctx context.Context) (string, error){
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // 3
	defer cancel()

	switch locale, err := locale_context(ctx);{
	case err!= nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func genFarewell__context(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	switch locale, err := locale_context(ctx);{
	case err!= nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func locale_context(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok{
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0{
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err() // 4
	case <-time.After(1* time.Minute):
		fmt.Println("terminate whole call-graph in locale time.After")
	}
	return "EN/US", nil
}



