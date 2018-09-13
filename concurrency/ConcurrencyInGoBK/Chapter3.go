package main

import (
	"sync"
	"bytes"
	"fmt"
	"time"
	"net/http"
	"math/rand"
)

//go 并发模式

func main(){
	//BytesBuffer()
	//SignalChildToStop()
	//ErrorHandling()
	//Pipeline()
	ChannelPattern()

}
func Count(ch chan int) {
	fmt.Println("Counting")
	ch <- 1
}

func BytesBuffer(){
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}

func SignalChildToStop(){
	// 通过一个只读的channel（done <-chan interface{}），父go routine线程通知子线程关闭相关通道或者终止某些操作

	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	v:=<-terminated
	fmt.Println("Done.")
	fmt.Println(v)
}

func OrChannel(channels ...<-chan interface{}) <-chan interface{} {
		//判断多个channel中是否有channel达到了关闭的条件（有值可取）


		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-OrChannel(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
}

func ErrorHandling(){
	type Result struct {
		Error error
		Response *http.Response
	}
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				fmt.Printf("after http.Get(%s) \n", url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					fmt.Println("done  -_-")
					return
				case results <- result:
					fmt.Println("result put into results channel")
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v \n", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func Pipeline() {
	//通过channel在多个通道间进行数据的传递
	// 通过done 变量在管道内进行级联关闭

	//generator 将一些离散的值映射到管道中， 方法：返回一个只读的channel，通过这个channel将数值依次输出
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values{
					select {
					case valueStream <- v:
					case <-done:
						return
					}
				}
			}
		}()
		return valueStream
	}
	take := func(done <- chan interface{}, number int, valueStream <- chan interface{}) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i:=0; i<number; i++{
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}
	repeatFn := func(done <- chan interface{}, fn func() interface{})<- chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string{
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream{
				select {
				case <-done:
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int, ) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			//time.Sleep(2*time.Second)
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i*multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int, ) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i+additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for i := range pipeline{
		//close(done)
		fmt.Printf("PPPP  %d\n", i)
	}
	for num := range take(done,10,  repeat(done, 1)){
		//close(done)
		fmt.Printf("%v  ", num)
	}
	rd := func() interface{}{
		return rand.Int()
	}
	for num := range  take(done, 10, repeatFn(done, rd)){
		fmt.Println(num)
	}
	var message string
	for token:= range toString(done, take(done,  5 ,repeat(done, "I", "am."))){
		message += token
	}
	fmt.Printf("message: %s...", message)
}

// 当完成Pipeline设计之后，Pipeline会在计算密集的stage部分发生block
// 可以利用Fan-Out，Fan-In 完成对Pipeline的不同的Stage的复用和并行化
// Fan-Out 利用多个goroutine 来处理来自pipeline的输入
//  //  //进行Fan-out的条件： stage不依赖stage之前的计算结果； stage需要花费长时间运行（order-independence、duration）
// Fan-In  将多个结果合并到一个channel中

func FanOutIn(){
//	这是一个关于FanOut的例子
//  这个例子中有两个stage：生成随机数、检查是否是素数
//	fanIn := func(
//		done <-chan interface{},
//		channels ...<-chan interface{},
//	) <-chan interface{} {
//		var wg sync.WaitGroup
//		multiplexedStream := make(chan interface{})
//		multiplex := func(c <-chan interface{}) {
//			defer wg.Done()
//			for i := range c {
//				select {
//				case <-done:
//					return
//				case multiplexedStream <- i:
//				}
//			}
//		}
//		// Select from all the channels
//		wg.Add(len(channels))
//		for _, c := range channels {
//			go multiplex(c)
//		}
//		// Wait for all the reads to complete
//		go func() {
//			wg.Wait()
//			close(multiplexedStream)
//		}()
//		return multiplexedStream
//		}
	}

func ChannelPattern(){
//	集中channel的使用的模式

	// 将channel c中的值 push到返回的只读的channel中
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values{
					select {
					case valueStream <- v:
					case <-done:
						return
					}
				}
			}
		}()
		return valueStream
	}
	take := func(done <- chan interface{}, number int, valueStream <- chan interface{}) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i:=0; i<number; i++{
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}

	// tee-channel
	tee := func(done <- chan interface{}, in <- chan interface{}) (_, _ <- chan interface{}) {
		//将一个channel的输入重定向到两个channel中
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in){
				// 本地拷贝 shadow
				var out1, out2 = out1, out2
				for i:=0; i<2; i++{
					select {
					case <- done:
					case out1<-val:
						out1 = nil
					case out2 <-val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)
	out1, out2 := tee(done, take(done,  4, repeat(done, 1, 2)))
	fmt.Println("tee channel 的模式下channel是如何传递数据的")
	for val1 := range out1{
		fmt.Printf("out1: %v, out2: %v \n", val1, <-out2)
	}

	fmt.Println("show how to bridge: ")
	// 将一个由多个channel组成的channel  解构成一个简单的channel (destructure the channel of channels into a simple channel)
	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		//valStream 返回所有的值
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				//从chanStream中拉取channel出来
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				// 使用从chanStream中拉出来的channel获取数据， 并将这些数据写入到valStream中
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	// bridge模式的应用实例 genVals产生具有十个channel的channel
	genVals := func() <-chan <-chan interface{}{
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i:=0; i<10; i++{
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v:= range bridge(nil, genVals()){
		fmt.Printf(" %v ", v)
	}

}

//	Queuing 用来将pipeline中的多个stage进行解耦
//  Queuing 并不会缩短整个流程的时间，但是减少多个流程之间相互因为耦合而相互堵塞的时间
func Queuing(){
	/*Queue的应用场景
		在Pipeline的入口处
		在Batch量可以提高性能的情况下
	*/

}


