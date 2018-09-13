package netchan

import (
	"fmt"
	"time"
)


func main(){
	msg := make(chan string)

	go Receive(msg)
	msg <- "hello"
	time.Sleep(time.Second*2)
	msg <- "ligang"
	time.Sleep(time.Second*1)
	msg <- "yes"
	close(msg)
	time.Sleep(time.Second*2)
}

func Receive(msgs chan string){
	for v := range msgs{
		fmt.Println(v)
	}
	fmt.Println("end")
}