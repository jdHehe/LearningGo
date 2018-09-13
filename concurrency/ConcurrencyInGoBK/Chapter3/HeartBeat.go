package main

import (
	"time"
	"fmt"
	"sync"
	"math/rand"
)

/*
通过channel实现心跳机制的模拟
time包实现定时事件，DoWork通过heartbeat  channel向调用者传递消息， 通过result  channel向调用者传递结果

*/

func DoWork(done <-chan interface{}, pulseInterval time.Duration) (<- chan interface{}, <- chan time.Time){
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2*pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct {}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-pulse:
					sendPulse()
				case results<-r:
					return
				}
			}
		}
		for i:=0; i<2; i++{
			select {
			case<-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func DoWorkWithNumber(done <- chan interface{}, nums ... int) (<-chan interface{}, <-chan int){
	heartbeat := make(chan interface{}, 1)
	initstream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(initstream)

		//time.Sleep(2*time.Second)

		for _, n := range nums{
			select {
			case heartbeat <- struct {}{}:
			default:
			}

			select {
			case <-done:
				return
			case initstream<-n:
			}
		}
	}()
	return heartbeat, initstream
}

func DoWorkWithReplicatedRequest(done<- chan interface{}, id int, wg *sync.WaitGroup, result chan <- int){
	//
	started := time.Now()
	defer wg.Done()

	// 模拟不同机器的延时不同这一属性
	simulatedLoadTime := time.Duration(1+rand.Intn(5))*time.Second
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <- done:
	case result <- id:
	}

	took := time.Since(started)
	if took < simulatedLoadTime{
		took = simulatedLoadTime
	}
	fmt.Printf(" %v took %v \n ", id, took)
}


func main(){
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done)
	})
	
	const timeout = 2 * time.Second
	heartbeat, results := DoWork(done, timeout/2)
	for{
		select {
		case _, ok := <-heartbeat:
			if ok == false{
				return 
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false{
				return
			}
			fmt.Printf("results %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}