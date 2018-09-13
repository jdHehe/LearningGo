package main


import (
	"sync"
	"fmt"
	"runtime"
	"testing"
	"time"
	"text/tabwriter"
	"os"
	"math"
	"bytes"
)

func main(){
	//GCDetail()
	//SyncMutex(5)
	//SyncRWMutex()
	//SyncCondBroadcast()
	//SyncOnce()
	//SyncPool()
	//ChannelBuf()
	//ChannelOwnerResponsibility()
	ChannelSelect()
}





func GoRoutineProcess(condition string){
	fmt.Println(condition)

	var wg sync.WaitGroup
	switch condition {
	case "withoutPara":
		// 这里goroutine极大可能会在for循环结束后才开始，也就是说，所有的goroutine拿到的location的值都是nanjing
		// 为什么for循环都结束了，goroutine还会拿得到for循环context中的变量？
		// go runtime is observant. 它将location的引用压到heap中，进而 goroutine 可以拿到location的值
		for _, location := range []string{"beijing", "chengdu", "nanjing"}{
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println(location)
			}()
		}
	case "withParaByValue":
		// 每次将 location的副本传入goroutine中所以correct
		for _, location := range []string{"beijing", "chengdu", "nanjing"}{
			wg.Add(1)
			go func(loca string) {
				defer wg.Done()
				fmt.Println(loca)
			}(location)
		}
	case "withParaByReference":
		// 传reference依然会出现第一个case的情况
		for _, location := range []string{"beijing", "chengdu", "nanjing"}{
			wg.Add(1)
			go func(loca *string) {
				defer wg.Done()
				fmt.Println(*loca)
			}(&location)
		}
	default:
		fmt.Println("wrong args")
	}
	wg.Wait()
}

func GCDetail(){
	// 有关garbage collect的memory 损耗问题

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }
	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})
	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}
	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}

func SyncMutex(interator int){
	var mux sync.Mutex
	count := 0

	increment := func() {
		mux.Lock()
		count++
		fmt.Printf("increment itreator %d\n", count)
		defer mux.Unlock()
	}

	decrement := func() {
		mux.Lock()
		count--
		fmt.Printf("decrement itreator %d\n", count)
		defer mux.Unlock()
	}

	var wg sync.WaitGroup
	wg.Add(interator)

	for i:=0; i<interator; i++{
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Add(interator)
	for i:=0; i<interator; i++{
		go func() {
			defer wg.Done()
			decrement()
		}()
	}
	wg.Wait()
}

func SyncRWMutex(){
	// 合理的应用RWMutex对读写进行区分  会在一定程度上加快整个性能
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i:=5; i>0; i--{
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}
	test := func(count int, mutex, rwMutex sync.Locker) time.Duration{
		var wg sync.WaitGroup
		wg.Add(count+1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i:=count; i>0; i--{
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 1,2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i:=0; i<20; i++{
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
			)
	}

}

func SyncCond(){
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 12)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		// 发送信号， 主线程中的wait得到释放
		c.Signal()

	}
	for i:=0; i<10; i++{
		c.L.Lock()
		for len(queue) == 2{
			//等待信号
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct {}{})
		go removeFromQueue(1*time.Second)
		c.L.Unlock()
	}

}

func SyncCondBroadcast(){
	// 利用 sync.Cond向所有注册的handler进行广播
	// 和Signal&Wait的模式类似，broadcast notify所有的wait的节点
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{ Clicked: sync.NewCond(&sync.Mutex{}) }
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}
	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})
	button.Clicked.Broadcast()
	clickRegistered.Wait()
}

func SyncOnce(){
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var onceA, onceB sync.Once
	onceA.Do(decrement)
	onceB.Do(increment)

	fmt.Printf("Count: %d\n", count)
}
func SyncPool(){
	/*注意事项
	• When instantiating sync.Pool, give it a New member variable that is thread-safe
	when called.
	• When you receive an instance from Get, make no assumptions regarding the
	state of the object you receive back.
	• Make sure to call Put when you’re finished with the object you pulled out of the
	pool. Otherwise, the Pool is useless. Usually this is done with defer.
	• Objects in the pool must be roughly uniform in makeup.
	*/
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorker  = 1024*1024
	var wg sync.WaitGroup
	wg.Add(numWorker)
	for i := numWorker; i>0; i--{
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)

			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

func ChannelBuf(){
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintf(&stdoutBuff, "Producer Done.")
		for i:=0; i<5; i++{
			fmt.Print("for  \t")
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
			time.Sleep(1*time.Second)
		}
		fmt.Print("\n")
	}()

	for integer := range intStream{
		fmt.Print("yyyy    \t")
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
	//fmt.Println("=-=-=============")
	//for i:=0; i<5;i++  {
	//	fmt.Println(<-intStream)
	//}

}

func ChannelOwnerResponsibility(){
	/* channel 的拥有者的责任
	1. Instantiate the channel.
	2. Perform writes, or pass ownership to another goroutine.
	3. Close the channel.
	4. Ecapsulate the previous three things in this list and expose them via a reader
	channel.
	 */

	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}
	resultStream := chanOwner()

	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}
func ChannelSelect() {
	// select 对channel的行为进行选择，可以接受消息也可以发送消息
	// 当 没有default的时候，select可能会一直等待对应的数据的到来
	// 设定一个timeout的时间也是一个常用的设计

	c1 := make(chan interface{}); close(c1)
	c2 := make(chan interface{}); close(c2)
	var c1Count, c2Count int
	for i := 10; i >= 0; i-- {
		select {
		case v:=<-c1:
			fmt.Print(v)
			c1Count++
		case <-c2:
			c2Count++
		case <-time.After(1*time.Second):
			fmt.Println("Time out")
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}