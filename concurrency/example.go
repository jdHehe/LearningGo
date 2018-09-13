package main
import (

	"sync"
	"fmt"

)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	//sem := make(chan int, 1)

	for i := 0; i < 3; i++ {

		go func(id int) {
			defer wg.Done()

			//sem <- 1

			for x := 0; x < 3; x++ {
				fmt.Println(id, x)
			}

			//<- sem
		}(i)
	}

	wg.Wait()
}