package driver

import (
	biu "github.com/jdHeHe/LearningGo/biubiu/context"
	"fmt"
	"flag"
)

/*
created by kakunka
2018/6/6
*/

/*
Driver use DAG to distribute works to clusters
Driver程序通过DAG将连续的任务划分为不同的stage，选取集群中worker 执行器来执行具体任务
*/

type BiuDriver struct {
	context *biu.BiuContext
	AllStageGroups []StageGroup
}





//  driver start working when the method Run() been called
//  run是整个driver程序运行的起始点
//  run 的整体运行为  针对context的task生成stage状态
//  由action触发整个stage的任务分发
//  driver等待action的结果，然后在进行下一次的任务
func (driver BiuDriver)Run(){
	var lastStageGroup  *StageGroup
	for _, stageGroup := range driver.AllStageGroups{
		if lastStageGroup != nil{
			for _, stage := range stageGroup.Stages{
				for i, shader := range stage.InputOfShard{
					shader.AgentLocaltion = lastStageGroup.Stages[i].Localtion
				}
			}
		}
		for _, stage := range stageGroup.Stages{
			agentLocation := driver.GetAgent()
			stage.Localtion = agentLocation
			for _, shader := range stage.OutputOfShard{
				shader.AgentLocaltion = agentLocation
			}
		}
	}
}

func (driver *BiuDriver)SetupAllStages(){
	stepGroups := driver.GroupStepsToStepGroups()
	driver.AllStageGroups = driver.FromStepGroupsToStages(stepGroups)
}




// 	将所有的steps分成一个个相互关联的Stages
//  每个Stage中
func (driver BiuDriver)GroupStepsToStepGroups() []StepGroup{
	fmt.Println("context中的steps数量  ", len(driver.context.Steps) )
	stepGroups := make([]StepGroup, 0)
	stepGroup  := NewStepGrup(len(stepGroups))
	for _, step := range driver.context.Steps{
		if !step.IsNarrow{
			fmt.Println("这里会新增stepGroup")
			if len(stepGroup.Steps) == 0{
				// 之前的steps
				stepGroup.AddStep(step)
				oldStepGroup := stepGroup
				stepGroups = append(stepGroups, *oldStepGroup)
				stepGroup = NewStepGrup(len(stepGroups))
				stepGroup.AddParent(oldStepGroup.ParentStepGroup, oldStepGroup)
				continue
			}
			oldStepGroup := stepGroup
			stepGroups = append(stepGroups, *oldStepGroup)
			stepGroup = NewStepGrup(len(stepGroups))
			stepGroup.AddParent(oldStepGroup.ParentStepGroup, oldStepGroup)
		}
		stepGroup.AddStep(step)
	}
	stepGroups = append(stepGroups, *stepGroup)
	fmt.Println("stepGroups的长度是", len(stepGroups))
	return  stepGroups
}
func (driver BiuDriver)FromStepGroupsToStages(stepGroups []StepGroup) []StageGroup{
	stageGroups := make([]StageGroup, 0)
	for _, stepGroup := range stepGroups{
		stages := make([]Stage, len(stepGroup.Steps[0].Tasks))
		stageGroup := NewStageGroup(len(stageGroups))
		for i, stage := range stages{
		// 设置Stage的输入
			if len(stepGroup.Steps[0].Tasks[0].InputShards) == 1{
				stage.InputOfShard = append(stage.InputOfShard, stepGroup.Steps[0].Tasks[0].InputShards[0])
			}else if len(stepGroup.Steps[0].Tasks[0].InputShards) > 1{
				// 从多个前驱 获得数据
				// 可能是 N:1  或  N:M
				for _, shader := range stepGroup.Steps[0].Tasks[0].InputShards{
					stage.InputOfShard = append(stage.InputOfShard, shader)
				}
			}else{
				fmt.Println("这个Stage ", stage.Id, " 没有输入数据")
			}
		//	设置stage的输出
			if (len(stepGroup.Steps[len(stepGroup.Steps)-1].Tasks[0].OutputShards) > 1){
			//	向多个stage输出
				for _, shader := range stepGroup.Steps[len(stepGroup.Steps)-1].Tasks[0].OutputShards {
					stage.OutputOfShard = append(stage.OutputOfShard, shader)
				}
			}else {
				// 只向一个Shader输出
				outShader := stepGroup.Steps[len(stepGroup.Steps)-1].Tasks[0].OutputShards[0]
				stage.OutputOfShard = append(stage.OutputOfShard, outShader)
			}
			for _, step := range stepGroup.Steps{
				stage.Tasks = append(stage.Tasks, *step.Tasks[i])
			}
			stageGroup.Stages = append(stageGroup.Stages, &stage)
		}
		stageGroups = append(stageGroups, *stageGroup)
	}
	return stageGroups
}

// 获取下一个能够接受任务的Agent
var (
	a = 0
)
func (driver BiuDriver)GetAgent() string{
	a++
	a %= 3
	agents := []string{":4000", ":3000", ":2000"}
	return agents[a]
}

func (driver BiuDriver)


func NewBiuDriver(context *biu.BiuContext) *BiuDriver{
	return &BiuDriver{
		context:context,
	}
}