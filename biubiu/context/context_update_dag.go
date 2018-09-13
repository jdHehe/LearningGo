package context

import (

)

//	这个文件中进行每一个步骤中task以及step相关的更新
//　针对宽、窄依赖进行不同的更新步骤
//  三种方式 1:1  1:n  n:1
//  1:n 发生在源dataset中，也就意味着在整个生命流程中一般只创建一次
//  在这些函数中主要是设置task的输入输出，通过将dataset的输入输出引流到task
//	这样当doTask中执行任务的时候，数据可以顺利的在dataset中流转


// 1:1窄依赖 每个Task只有一个输入和一个输出
func (context *BiuContext)OneInputForOneOutput(input * DataSet, output * DataSet, step * Step)  {
	if output != nil{
		output.Step = step
		step.Outputs = append(step.Outputs, output)
	}
	if input != nil{
		step.Inputs = append(step.Inputs, input)
	}

	// 将input和output的shard都绑定到tasks中
	if input != nil {
		for i, shard := range input.Datas {
			task := step.NewTask()
			task.InputShards = append(task.InputShards, shard)
			task.OutputShards = append(task.OutputShards, output.Datas[i])
		}
	}
}

// 1:n Task一个输入，多个输出
func (context *BiuContext)FromOneDsToMany(input * DataSet, output * DataSet, step * Step){
	if output != nil{
		output.Step = step
		step.Outputs = append(step.Outputs, output)
	}
	if input != nil{
		step.Inputs = append(step.Inputs, input)
	}
	task := step.NewTask()
	if input != nil{
		task.InputShards = append(task.InputShards, input.Datas[0])
	}
	for _, shader := range  output.Datas{
		task.OutputShards = append(task.OutputShards, shader)
	}

}

//n:1 Task多个输入，一个输出
func (context *BiuContext)FromManyDsToOne(input * DataSet, output * DataSet, step * Step){
	if output != nil{
		output.Step = step
		step.Outputs = append(step.Outputs, output)
	}
	if input != nil{
		step.Inputs = append(step.Inputs, input)
	}
	task := step.NewTask()
	if output != nil{
		task.OutputShards = append(task.OutputShards, output.Datas[0])
	}
	for _, shader := range  input.Datas{
		task.InputShards = append(task.InputShards, shader)
	}
}