package mapreduce

import (
	"fmt"
	"strconv"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	info := <-mr.registerChannel
	fmt.Println(" mr.registerChannel 的内容是  ", info)
	fmt.Println(mr.workers)
	numWorkers := len(mr.workers) //worker的个数
	fmt.Println("注册的worker的数目  ", numWorkers)
	if numWorkers == 0{
		panic(" no useful worker ")
	}


	//存储所有的调用结果，鉴于rpc call采用： 返回值true or false标志调用是否正确完成
	// 需要将所有的call的调用结果都存储起来
	rpcResult := make(chan bool, ntasks)

	defer close(rpcResult)
	fmt.Println("这次的任务的总数 ", ntasks)
	for i:=0; i<ntasks; i++{
		go func(index int) {
			for {
				// 一些配置项     告诉mapper和reducer如何 输出/输入 文件
				//worker的编号
				workerNumber := index%len(mr.workers)
				workerName := mr.workers[workerNumber]
				args := new(DoTaskArgs)
				args.Phase = phase
				args.JobName = mr.jobName
				args.NumOtherPhase = nios
				args.File = mr.files[index]
				args.TaskNumber = index
				fmt.Println("进行tasker任务的具体执行 第", index, "个任务", " 任务了类型", phase, "文件", args.File)
				result := call(workerName, "Worker.DoTask", args, new(struct{}))
				if result == true {
					rpcResult <- result
					return
				} else {
					//	如果这个worker失败，需要向master 报告这个错误，将错误的worker移除，并将分配给他的任务分配给其他worker
					//  防止重复的移除，我们需要先判断这个worker是否还存在于master中
					if  workerNumber<len(mr.workers) && mr.workers[workerNumber] == workerName {
						if len(mr.workers) == 1 {
							panic("最后一个worker也down了没有worker可以使用了" + mr.workers[0])
						}
						mr.Mutex.Lock()
						mr.workers = append(mr.workers[:workerNumber], mr.workers[workerNumber+1:]...)
						mr.Mutex.Unlock()
						fmt.Println("worker shutdown shift another one ===================")
					}
				}
			}
		}(i)
	}
	for i:=0; i<ntasks; i++{
		v := <-rpcResult
		if !v {
			panic(strconv.Itoa(i)+" network still call fails ")
		}
	}

	// 换个思路  创建channel切片，通过切片的下标对应任务的下标，便于对失败的任务进行回溯
	//rpcResult := make([]chan bool, ntasks)
	//successNumber := 0
	//for i:=0; i<ntasks; i++{
	//	numWorkers = len(mr.workers)
	//		go func(index int) {
	//			// 一些配置项     告诉mapper和reducer如何 输出/输入 文件
	//			args := new(DoTaskArgs)
	//			args.Phase = phase
	//			args.JobName = mr.jobName
	//			args.NumOtherPhase = nios
	//			args.File = mr.files[index]
	//			args.TaskNumber = index
	//			fmt.Println("进行tasker任务的具体执行 第", index, "个任务", " 任务了类型", phase, "文件", args.File)
	//			result := call(mr.workers[i%numWorkers], "Worker.DoTask", args, new(struct {}))
	//			if result == true{
	//				ch <- result
	//			}else{
	//			//	如果这个worker失败，需要向master 报告这个错误，将错误的worker移除，并将分配给他的任务分配给其他worker
	//
	//			}
	//		}(i)
	//}


	debug("Schedule: %v phase done\n", phase)
}
