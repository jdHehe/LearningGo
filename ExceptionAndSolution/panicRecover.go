
package main

import (
	"fmt"
	"sync"
	"time"
)
func main() {
	var wg sync.WaitGroup
	i := 0

	wg.Add(2)
	go f(&wg, &i)
	go e(&wg, &i)
	wg.Wait()
}


func f(wg *sync.WaitGroup, i * int) {
	defer func(wg *sync.WaitGroup) {     //必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err)    //这里的err其实就是panic传入的内容，55
		}
		*i++
		fmt.Println("d")
		wg.Done()
	}(wg)
	*i++
	fmt.Println("a")
	panic(55)
	*i++
	fmt.Println("b")
	fmt.Println("f")

}

func e(wg *sync.WaitGroup, i *int){
	defer wg.Done()
	time.Sleep(1*time.Second)
	fmt.Println("_++))+()+")
	*i++
}