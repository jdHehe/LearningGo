package context

import (
	"sync"
	"github.com/jdHeHe/LearningGo/biubiu/util"
	"os"
	"bufio"
	"io"
	"reflect"
	"fmt"
)

type BiuContext struct {
	sync.Mutex
	DataSets 	[]*DataSet 	// 这个任务中的所有的DataSet集合
	jobName 	string				// 运行的作业的名称
	done 		chan struct{}		// 终止作业
	Steps   	[]*  Step  // 这次作业的整个的DAG的各个环节
}

//	从本地文件构建一个DataSet
//  TextFile按行读取文本文件
func (context *BiuContext)TextFile(fileName string, shards int) *DataSet{
	//	定义一个通过channel传递消息的回调函数
	// 	fn按行读取文件并传递给out channel
	fn := func(out chan string) {
		file, err := os.Open(fileName)
		if err != nil{
			panic(err)
		}
		bfRd := bufio.NewReader(file)
		for{
			line, err := bfRd.ReadString('\n')
			if err != nil || err == io.EOF{
				break
			}
			//fmt.Println("从文件读取一行数据", line)
			out<-line
		}
		fmt.Println("文件读取结束")
	}
	ds := context.Source(fn, shards)
	return ds
}

// 根据fn函数定义的数据获取方式生成一个DataSet
func (context *BiuContext)Source(fn interface{}, shards int) (ds*DataSet){
	type_ := util.FuncType(fn)
	step := context.NewStep()
	step.Name = "source"
	//  一对多
	step.IsNarrow = false
	ds = NewDataSet(context, type_)
	ds.SetupShard(2)
	ds.Step = step
	context.FromOneDsToMany(nil, ds, step)
	step.Function = func(task *  Task) {
		chanType := reflect.ChanOf(reflect.BothDir, type_)
		outChan  := reflect.MakeChan(chanType, 0)
		funcEle  := reflect.ValueOf(fn)

		fmt.Println("chanType:",chanType, " outchan:",reflect.TypeOf(outChan), " funcEle:",funcEle)
		//  通过outchan 从fn中获取数据，然后将数据传递给task的OutPutChan
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer outChan.Close()
			funcEle.Call([]reflect.Value{outChan})
		}()

		wg.Add(1)
		go func() {
			defer  wg.Done()
			defer func() {
				for _,sharder := range task.OutputShards{
					sharder.Close()
				}
			}()
			i := 0
			for {
				if v, OK := outChan.Recv(); OK{
					fmt.Println(v)
					task.OutputShards[i].WriteChan.Send(v)
					i++
					if i== shards{
						i=0
					}
				}else {
					fmt.Println("outchan recv false !!!!!!!!")
					break
				}
			}
		}()
		wg.Wait()
		fmt.Println("source ", "step:", task.Step.Name, "taskId:", task.Id)
	}
	return
}


func NewBiu(string string)  *BiuContext{
	return &BiuContext{
		jobName:string,
		done:make(chan struct{}),
	}
}
func (context *BiuContext)NewStep() *Step{
	step := &  Step{
		Id:len(context.Steps),
		IsNarrow:true,
	}
	context.Steps = append(context.Steps, step)
	return step
}
func (context *BiuContext)Run(){
	fmt.Println("context.steps 的长度:", len(context.Steps))
	for i, step := range context.Steps{
		fmt.Println("step" , i, "的task长度", len(step.Tasks))
		for j,task :=  range step.Tasks{
			fmt.Println("step ", i," task",j)
			go task.RunTask()
		}
	}
}