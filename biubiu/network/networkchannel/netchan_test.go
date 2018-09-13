package netchan

import (
	"testing"
	"fmt"
	"os"
	"reflect"
	"sync"
	"github.com/jdHeHe/LearningGo/biubiu/context"
)

func TestSenderAndReceiver(t *testing.T){
	kv := context.KeyValue{
		Key:"Key",
		Value:"Value",
	}
	receiver, err := GetReceiverChannelByReflectValue("localhost:4000", kv)
	if err != nil{
		fmt.Println("GetSenderChannelByReflectValue     ", err)
		os.Exit(1)
	}
	sender, err := GetSenderChannelByReflectValue(":4000", kv)
	var wg sync.WaitGroup
	wg.Add(1)
	DoChan(sender, receiver, &wg)
	sender2, _ := GetSenderChannelByReflectValue(":4000", kv)
	wg.Add(1)
	DoChan(sender2, receiver, &wg)
	wg.Wait()
}

func DoChan(Sender chan<- interface{}, Receiver <-chan interface{}, wg *sync.WaitGroup){
	defer wg.Done()
	kv := context.KeyValue{
		Key:"Key1",
		Value:"Value1",
	}
	fmt.Println("DoChan===================    ", kv)
	send := reflect.ValueOf(Sender)
	fmt.Println(send.String())
	kv_value := reflect.ValueOf(kv)
	fmt.Println(kv_value.String())
	send.Send(kv_value)
	//Sender <-
	value := reflect.ValueOf(Receiver)
		//<- Receiver
	if v, OK :=value.Recv();OK{
		fmt.Println("value is  ")
		fmt.Println(v.Interface().(context.KeyValue).Key)
	}

}