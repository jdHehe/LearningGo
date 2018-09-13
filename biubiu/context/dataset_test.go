package context

import (
	"testing"
	"strings"
	"fmt"
	"sync"
	"github.com/jdHeHe/LearningGo/biubiu/util"
	"time"
)


func TestMap(t *testing.T){
	biu := NewBiu("biu")
	biu.TextFile("D://pg-huckleberry_finn.txt", 2).Map(	func(line string, ch chan KeyValue) {
		fmt.Println("in map function line is", line)
		for _, token := range strings.Split(line, " ") {
			ch <- KeyValue{
				Key:token,
				Value:1,
			}
		}
	}).LocalSort(nil).LocalReduceByKey(func(x, y int) int{
		//fmt.Println("call LocalReduce  ", x, y)
		return x+y
	})
	go biu.Run()
	var wg sync.WaitGroup
	for i, shard := range biu.DataSets[len(biu.DataSets)-1].Datas{
		wg.Add(1)
		fmt.Println("----------------------result------------------")
		go func(shard *DataShard, index int) {
			defer wg.Done()
			for {
				if v, OK := shard.WriteChan.Recv(); OK {
					//fmt.Println(v.String())
					fmt.Println( index,"  v:   key:", v.Interface().(KeyValue).Key, "  value:", v.Interface().(KeyValue).Value)
				} else {
					fmt.Println("============++++++++++++++++ break")
					break
				}
			}
			fmt.Println("----------------------result------------------")
		}(shard, i)
	}
	wg.Wait()

}
func TestAddOutput(t *testing.T){
	biu := NewBiu("biu")
	biu.TextFile("D://hello.txt", 2).Map(	func(line string, ch chan string) {
		fmt.Println("in map function line is", line)
		for _, token := range strings.Split(line, ":") {
			fmt.Print(" 把 ", token, "塞到channel中")
			ch <- token
		}
	}).Map(func(key string, ch chan int) {
		fmt.Print("key: ", key, "  value:", 1)
		ch <- 1
	})
	var wg  sync.WaitGroup
	//biu.TextFile("D://hello.txt", 2)
	go biu.Run()
	result := 0
	for _, shard := range biu.DataSets[len(biu.DataSets)-1].Datas{
		wg.Add(1)
		fmt.Println("----------------------result------------------")
		go func(shard *DataShard) {
			defer wg.Done()
			for {
				if v, OK := shard.WriteChan.Recv(); OK {
					//fmt.Println(v.String())
					result += v.Interface().(int)
				} else {
					fmt.Println("============++++++++++++++++ break")
					break
				}
			}
			fmt.Println("----------------------result------------------")
		}(shard)
	}

	wg.Wait()
	fmt.Println(result)
}

func TestKVType(t *testing.T){
	type KeyValue struct {
		Key string
		Value int
	}
	biu := NewBiu("biu")
	biu.TextFile("D://hello.txt", 2).Map(	func(line string, ch chan KeyValue) {
		for _, token := range strings.Split(line, ":") {
			ch <- KeyValue{
				Key:token,
				Value:1,
			}
		}
	}).Map(func(key KeyValue, ch chan int) {
		fmt.Print("key: ", key, "  value:", 1)
		ch <- 1
	})
	biu.Run()
}
func  TestLocalSort(t *testing.T){

	//var result chan KeyValue
	biu := NewBiu("biu")
	//pg-huckleberry_finn
	biu.TextFile("D://pg-les_miserables.txt", 2).Map(	func(line string, ch chan KeyValue) {
		//fmt.Println("map line ",line)
		for _, token := range strings.Split(line, " ") {
			ch <- KeyValue{
				Key:token,
				Value:1,
			}
		}
	}).LocalSort(nil).LocalReduceByKey(func(x, y int) int{
		//fmt.Println("call LocalReduce  ", x, y)
		return x+y
	}).MergeSorted(func(x, y int) int{
		fmt.Println("call MergeSort  ", x, y)
		return x+y
	})
	//.AddOutput(result)

	fmt.Println("++++++++++++++")
	for _,ds := range  biu.DataSets{
		fmt.Println(ds.Type)
	}
	fmt.Println("biu.Run========================")
	filePath := "D://serilization.gob"
	err := util.Serilization(filePath, biu)
	if err != nil{
		fmt.Println("serilization  ", err)
		return
	}
	time.Sleep(1*time.Second)
	var newBiu BiuContext
	err = util.DeSerilization(filePath, newBiu)
	if err != nil{
		fmt.Println("deserilization ", err)
		return
	}
	fmt.Println(newBiu.jobName)


	//go biu.Run()
	//var wg sync.WaitGroup
	//for i, shard := range biu.DataSets[len(biu.DataSets)-1].Datas{
	//	wg.Add(1)
	//	fmt.Println("----------------------result------------------")
	//	go func(shard *DataShard, index int) {
	//		defer wg.Done()
	//		for {
	//			if v, OK := shard.WriteChan.Recv(); OK {
	//				//fmt.Println(v.String())
	//				fmt.Println( index,"  v:   key:", v.Interface().(KeyValue).Key, "  value:", v.Interface().(KeyValue).Value)
	//			} else {
	//				fmt.Println("============++++++++++++++++ break")
	//				break
	//			}
	//		}
	//		fmt.Println("----------------------result------------------")
	//	}(shard, i)
	//}
	//wg.Wait()


}

func TestStepSerilize(t *testing.T){


}