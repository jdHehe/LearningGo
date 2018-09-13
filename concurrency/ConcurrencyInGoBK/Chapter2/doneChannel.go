package main

import (
	"sync"
	"fmt"
	"time"
)
// 通过  同一个channel （done channel）来控制整个call-graph的流向


func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil{
			fmt.Printf(" %v ", err)
			return
		}
	}()

	wg.Wait()
}
func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil{
		return err
	}
	fmt.Printf(" %s world ! \n", farewell)
	return nil
}
func printGreeting(done <- chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil{
		return err
	}
	fmt.Printf(" %s world ! \n", greeting)
	return nil
}
func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done);{
	case err!= nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done);{
	case err!= nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(3 * time.Second):
		fmt.Println("terminate whole call-graph in locale time.After")
	}
	return "EN/US", nil
}

