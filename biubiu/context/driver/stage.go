package driver



import (
    context "github.com/jdHeHe/LearningGo/biubiu/context"
	"reflect"
	"fmt"
	"github.com/jdHeHe/LearningGo/biubiu/network/networkchannel"
)

/*
created by kakunka
2018/6/6
*/

// stage.go 完成works的组装
type Stage struct {
	// InputOfShard OutputOfShard 是一个network channel，
	// 分别从其他Agent接受数据和向下一个 Agent发送数据
	InputOfShard 	[]*context.DataShard
	OutputOfShard 	[]*context.DataShard
	Tasks	[] context.Task
	Id int
	InputType 	interface{}
	OutputType 	interface{}
	Localtion 	string
}

// 设置InputOfShard
//func (stage *Stage)SetInput(target string, Instance interface{}) error{
//	receiver, err := netchan.GetReceiverChannelByReflectValue(target, Instance)
//	if err != nil{
//		return err
//	}
//	stage.InputOfShard = append(stage.InputOfShard, reflect.ValueOf(receiver))
//	return nil
//}
//func (stage *Stage)SetOutput(target string, Instance interface{}) error{
//	sender, err := netchan.GetSenderChannelByReflectValue(target, Instance)
//	if err != nil{
//		return err
//	}
//	stage.OutputOfShard = append(stage.OutputOfShard, reflect.ValueOf(sender))
//	return nil
//}
func (stage *Stage)UpdateInputOutputShard() error{
	for _, shader := range stage.InputOfShard{
		receiver, err := netchan.GetReceiverChannelByReflectValue(shader.AgentLocaltion, stage.InputType)
		if err != nil{
			return err
		}
		shader.WriteChan = reflect.ValueOf(receiver)
	}

	for _,shader := range stage.OutputOfShard{
		sender, err := netchan.GetSenderChannelByReflectValue(shader.AgentLocaltion, stage.OutputType)
		if err != nil{
			return err
		}
		shader.WriteChan = reflect.ValueOf(sender)
	}
	return nil
}



type StageGroup struct {
	Stages []*Stage
	Id int
	//InputLocaltion string
	//OutputLocaltion string
}

func (stageGroup *StageGroup)SetInput(){

}

func NewStageGroup(id int) *StageGroup{
	return &StageGroup{
		Id:id,
		Stages:make([]*Stage, 0),
	}
}


func NewStage(id int)*Stage{
	return &Stage{
		InputOfShard:make([]*context.DataShard, 0),
		OutputOfShard:make([]*context.DataShard, 0),
		Tasks:make([]context.Task, 0),
		Id:id,
	}
}

type StepGroup struct {
	Steps []*context.Step
	ParentStepGroup []*StepGroup
	Id int
}

func (stepGroup *StepGroup)AddStep(step *context.Step){
	fmt.Println("Append step to StepGroup")
	stepGroup.Steps = append(stepGroup.Steps, step)
	fmt.Println("after append steps长度:" , len(stepGroup.Steps))
}
func (stepGroup *StepGroup)AddParent(grandparent []*StepGroup, stepGroup_ *StepGroup){
	stepGroup.ParentStepGroup = append(stepGroup.ParentStepGroup, grandparent[:]...)
	stepGroup.ParentStepGroup = append(stepGroup.ParentStepGroup, stepGroup_)
}

func NewStepGrup(id int) *StepGroup{
	return &StepGroup{
		Steps:make([]*context.Step, 0),
		ParentStepGroup:make([]*StepGroup, 0),
		Id:id,
	}
}

func GroupWorksToStage(){

}
