package driver

import (
	"testing"
	"strings"
	"github.com/jdHeHe/LearningGo/biubiu/context"
	"fmt"
)

func TestGroupWorksToStage(t *testing.T) {
	biu := context.NewBiu("biu")
	//pg-huckleberry_finn
	biu.TextFile("D://pg-les_miserables.txt", 2).Map(	func(line string, ch chan context.KeyValue) {
		//fmt.Println("map line ",line)
		for _, token := range strings.Split(line, " ") {
			ch <-context. KeyValue{
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
	for _, step := range biu.Steps{
		for _, task := range step.Tasks{
			fmt.Println(task.Id, task.Step.Name)
		}
	}



	driver := NewBiuDriver(biu)
	fmt.Println("NewDriver")
	groups := driver.GroupStepsToStepGroups()
	for _, steps := range groups{
		fmt.Println(" stepsID ", steps.Id, "  ", "steps 的长度: ", len(steps.Steps))
		for _, step := range steps.Steps{
			fmt.Print("  stepId:",step.Id,"  stepName:", step.Name)
			for _, task := range step.Tasks{
				fmt.Print(" taskId:", task.Id, "task.StepName:",task.Step.Name,"         ")
			}
			fmt.Println()
		}
	}

	stageGroups := driver.FromStepGroupsToStages(groups)
	for _, stageGroup := range  stageGroups{
		fmt.Println(len(stageGroup.Stages))
	}
}

func TestDriverRun(t *testing.T){
	biu := context.NewBiu("biu")
	//pg-huckleberry_finn
	biu.TextFile("D://pg-les_miserables.txt", 2).Map(	func(line string, ch chan context.KeyValue) {
		//fmt.Println("map line ",line)
		for _, token := range strings.Split(line, " ") {
			ch <-context. KeyValue{
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
	driver := NewBiuDriver(biu)
	driver.SetupAllStages()
	driver.Run()
	for _, stageGroup := range driver.AllStageGroups{
		for _, stage := range stageGroup.Stages{
			fmt.Println("stage input  len ", len(stage.InputOfShard))
			for _,shader := range stage.InputOfShard{
				fmt.Print("  input localtion ", shader.AgentLocaltion)
			}
			fmt.Print("  Id:", stage.Id, "  localtion:", stage.Localtion, " tasks ", len(stage.Tasks))
		}
		fmt.Println()
	}


}

