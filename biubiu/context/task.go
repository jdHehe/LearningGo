package context

import (
	"sync"
)

//	Task执行的具体的任务
type Task struct {
	Id  		 	int
	//InputChan 		[]chan reflect.Value	// 输入【暂且以文件名代替】、
	//OutputChan		reflect.Value			// 输出
	Step 			*Step					// Step定义了Task具体要执行的操作
	// Task连接两个DataSet，InputShards是上一个DataSet的数据分片，OutputShard是下一个DataSet数据分片
	InputShards 	[]* DataShard
	OutputShards	[]* DataShard
	sync.Mutex
}
func (task *Task)RunTask(){
	// 执行step 的操作
	task.Step.Function(task)
}





// 每个DataSet的Transform和action都会产生一个新的Step
type Step struct {
	Name 		string
	Id   		int
	Tasks 		[]*Task 	//Tasks 存储的是这个step包括的Task
	Function 	func(task *Task)
	Outputs 	[]* DataSet //这个step的输入
	Inputs      []* DataSet //这个step的输出
	IsNarrow	bool
}

func (s *Step)RunStep(){
	var wg sync.WaitGroup
	for i, t := range s.Tasks {
		wg.Add(1)
		go func(i int, t *Task) {
			defer wg.Done()
			t.RunTask()
		}(i, t)
	}
	wg.Wait()
}
func (s *Step) NewTask() *Task{
	task := new(Task)
	task.Step = s
	s.Tasks = append(s.Tasks, task)
	task.Id = len(s.Tasks)
	return task
}