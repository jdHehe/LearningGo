package main

import (
	//context2 "github.com/jdHeHe/LearningGo/biubiu/context"
	"fmt"
	"os"
)

func main(){
	//biu := context2.NewBiu("biu")
	//biu.TextFile("abc", 3).Map(func(line string, ch chan string){})
	////.Reduce(func() {})
	//fmt.Println(biu.Steps)
	//
	//fmt.Println("==========================")
	//for _, step := range biu.Steps{
	//	fmt.Println()
	//	for _, task := range step.Tasks{
	//		//fmt.Print(task.Id, task.Step.Name)
	//		fmt.Println("task.InputChan len ", len(task.InputChan), "task.OutputChan len")
	//		task.RunTask()
	//	}
	//}
	//
	//for _, dataSet := range  biu.DataSets{
	//	fmt.Println("------shard的长度------", len(dataSet.Datas), "   Id:", dataSet.Id, " step:", dataSet.Step.Name)
	//	fmt.Println("------------------")
	//}
	//fmt.Println("=========================")

	fmt.Println(os.Args[0])
}
