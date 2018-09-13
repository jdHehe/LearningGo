package context

import (
	"reflect"
	"github.com/jdHeHe/LearningGo/biubiu/util"
	"fmt"
	"sync"
	"sort"
)

type DataSet struct {
	Id 					int						//编号
	Step	 			*Step
	parentDataset 		[]*DataSet
	context 			*BiuContext		//作业的上下文文件
	Type				reflect.Type			//生成这个DataSet的函数的输出类型
	Datas				[]*DataShard
	sync.Mutex
	Sorted				bool
}

/*
Map
假定f的形式为f(string, chan string)
Map 对应的DataSet的类型就是f函数的返回值的类型
（如果函数以channel的参数的形式获得返回值则  DataSet的类型为channel的数据类型）
*/
func (dataset *DataSet)Map(f interface{})  *DataSet{
	type_ := util.FuncType(f)
	//fmt.Println("mapppppppppp", type_)
	newDataset := NewDataSet(dataset.context, type_)
	newDataset.SetupShard(len(dataset.Datas))
	step := dataset.context.NewStep()
	step.Name = "Map"
	dataset.context.OneInputForOneOutput(dataset, newDataset, step)
	step.Function = func(task *Task) {
	//	map处理的数据 是一对一的
	//	进行map的逻辑 从task取数据，调用f函数，将结果写到写一个shard
		defer func() {
			fmt.Println("shader.close============================")
			for _,shader := range task.OutputShards{
				shader.Close()
			}
		}()
		fmt.Println("=============================task.InputShader  ",len(task.InputShards))
		fmt.Println(task.Id)
		var wg  sync.WaitGroup
		fmt.Println("wg:   ", wg)
		//for i, shard := range task.InputShards{
		//	// 实际上在map这一步的操作中，该task的InputsShard应当只有一个
		//	// 也就是说这个for loop应当执行一次
			shard := task.InputShards[0]
			invokeMapFunc := _MapperFunction(f, &task.OutputShards[0].WriteChan)
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, shard_ *DataShard) {
				defer wg_.Done()
				for{
					if v, OK := shard_.WriteChan.Recv(); OK{
						fmt.Println("<<<<<<<<<<<<<<<<<<<<<   ", v)
						invokeMapFunc(v)
					}else{
						break
					}
				}
			}(&wg, shard)
		//}
		wg.Wait()
	}
	return newDataset
}
func _MapperFunction(f interface{}, outChan *reflect.Value) func(input  reflect.Value) {
	return func(input  reflect.Value) {
		//fmt.Println("_mapperFucntion")
		//fmt.Println(input.String())
		fn := reflect.ValueOf(f)
		switch input.Type(){
		case KeyValueType:
			fmt.Println("KeyValueType call function")
			kv := input.Interface().(KeyValue)
			//fn.Call([]reflect.Value{reflect.ValueOf(kv.Key), reflect.ValueOf(kv.Value), *outChan})
			fn.Call([]reflect.Value{reflect.ValueOf(kv),  *outChan})
		default:
			//fmt.Println("default function call")
			fn.Call([]reflect.Value{input, *outChan})
		}




		}
}

// ReduceByKey 接受的DataSet的类型是KV的形式的
// ReduceByKey 返回的亦是KV的形式的DataSet
func (dataset *DataSet)ReduceByKey(f interface{}) *DataSet{
	return dataset.LocalSort(nil).MergeReduce(nil)
}


func (dataset *DataSet)Reduce(f interface{}) *DataSet{
	//return dataset.LocalReduce(f).MergeReduce(f)
	return nil
}

//	在本地执行Reduce操作
func (dataset *DataSet) LocalReduceByKey(f interface{}) *DataSet{
	newDataset := NewDataSet(dataset.context, dataset.Type)
	newDataset.SetupShard(len(dataset.Datas))
	step := dataset.context.NewStep()
	step.Name = "LocalReduceByKey"
	dataset.context.OneInputForOneOutput(dataset, newDataset, step)
	step.Function = func(task *Task) {
		fmt.Println("localreducebykey")
		outChan := task.OutputShards[0].WriteChan
		defer outChan.Close()
		var wg sync.WaitGroup
		fn := reflect.ValueOf(f)
		for _, shader := range task.InputShards{
			wg.Add(1)
			flagExist := false //标志当前的key是否和前一个key相同
			var previouskey interface{}
			var tmpResult reflect.Value
			go func(flag_exist bool, previous_key interface{}, tmp_result reflect.Value) {
				_handleShaderWriteChannelReceive(&wg, func(v reflect.Value) {
					//	 将相同的key的value通过用户定义的组合方式fold在一起，然后仍然以kv的形式发给下一个DataSet
					//fmt.Println("handleV LocalReduce  ",v)
					// 暂存合并结果
					kv := v.Interface().(KeyValue) //从上一个Datatset传入的DataSet需要时kv的形式
					if flag_exist {
						//fmt.Println("KV: ",  kv)
						if (reflect.DeepEqual(previous_key, kv.Key)){
							//	当前key和前一个key相同
							//fmt.Println("deepEqual=======")
							//fmt.Println("key相同  ",reflect.TypeOf(previous_key), reflect.TypeOf(kv.Key))
							//fmt.Println("key相同  ",reflect.ValueOf(tmp_result), reflect.ValueOf(kv.Key))
							tmp_result = fn.Call([]reflect.Value{tmp_result, reflect.ValueOf(kv.Value)})[0]
						}else {
							//fmt.Println("key不相同  ",reflect.TypeOf(previous_key), reflect.TypeOf(kv.Key))
							//fmt.Println("key不相同  ",reflect.ValueOf(previous_key), reflect.ValueOf(kv.Key))
							//fmt.Println("outChan.send:  ", previous_key, tmp_result)
							outChan.Send(reflect.ValueOf(KeyValue{Key:previous_key, Value:tmp_result.Interface()}))
							previous_key = kv.Key
							tmp_result = reflect.ValueOf(kv.Value)
						}
					}else {
						flag_exist = true
						previous_key = kv.Key
						tmp_result = reflect.ValueOf(kv.Value)
					}
				}, func() {
					// 所有数据接收结束，将最后一个数据或一组数据发送出去
					if flag_exist{
						outChan.Send(reflect.ValueOf(KeyValue{Key:previous_key, Value:tmp_result.Interface()}))
					}
				},shader)

			}(flagExist, previouskey, tmpResult)
		}
		wg.Wait()
	}
	return newDataset
}

// 	从多个sharder获得数据， 然后合并到新的dataset
func (dataset *DataSet) MergeReduce(f interface{}) *DataSet {
	newDataset := NewDataSet(dataset.context, dataset.Type)
	newDataset.SetupShard(1)
	step := dataset.context.NewStep()
	step.Name = "MergeReduce"
	dataset.context.FromManyDsToOne(dataset, newDataset, step)
	step.Function = func(task *Task) {
	//	MergeReduce对应的Task只有一个输出shard，所以只需要取第一个shader的WriteChan
		outChan := task.OutputShards[0].WriteChan
		defer outChan.Close()
		var wg sync.WaitGroup
		for _,shader := range task.InputShards{
			wg.Add(1)
			go _handleShaderWriteChannelReceive(&wg, func(v reflect.Value) {
			//	 接受前一个Dataset的shader传来的数据，经过处理，发送到下一个shader
			//	kv :=

			}, nil,shader)
		}
		wg.Wait()

	}
	return newDataset
}

// 对本地的经过sort的DataSet进行Merge操作
// 输入为具有多个Shard的
// 暂时采用全局的mapping 来存储各个Shard的KV合并的结果
func (dataset *DataSet)MergeSorted(mergeFn interface{}) *DataSet{
	newDataset := NewDataSet(dataset.context, dataset.Type)
	newDataset.SetupShard(1)
	step :=  dataset.context.NewStep()
	step.Name = "MergeSort"
	// 多对一
	step.IsNarrow = false
	dataset.context.FromManyDsToOne(dataset, newDataset, step)
	step.Function = func(task *Task) {
		//将多个输入shader的数据进行merge
		var wg sync.WaitGroup
		var lock sync.RWMutex
		mergeRes := make(map[string]interface{})
		fn := reflect.ValueOf(mergeFn)

		outChan := task.OutputShards[0].WriteChan
		defer outChan.Close()
		fmt.Println("+++++++++++++++++++++++++++++开始merge   MergeSort task.InputShards ", len(task.InputShards))
		for _, shader := range task.InputShards{
			wg.Add(1)
			go func(merge_res map[string]interface{}, shader_ *DataShard) {
				defer fmt.Println("结束一个merge routine")
				_handleShaderWriteChannelReceive(&wg, func(v reflect.Value) {
					kv := v.Interface().(KeyValue)
					fmt.Println("Key.........", kv.Key.(string))
					// 采用全局锁，极其的不高效
					lock.Lock()
					if v, OK := (merge_res)[kv.Key.(string)]; OK{
						fmt.Println("Map中存在 item key: ", kv.Key , v, "vtype ", reflect.TypeOf(v), " type of kv.Value ", reflect.TypeOf(kv.Value), kv.Value)
						res := fn.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(kv.Value)})[0]
						(merge_res)[kv.Key.(string)] = res.Interface()
					}else{
						fmt.Println("Map中第一次出现 item  key",  kv.Key," kv.value 的类型", reflect.TypeOf(kv.Value))
						(merge_res)[kv.Key.(string)] = kv.Value
					}
					lock.Unlock()
				}, nil, shader_)
			}(mergeRes, shader)


		}
		fmt.Println("wg.wait....................................................")
		wg.Wait()
		fmt.Println("after wg.wait................................................")
		fmt.Println("mergeRes--------------------------------  ", mergeRes)


		for k,v := range mergeRes {
			fmt.Println("outChan.send======", k, "  ", v)
			outChan.Send(reflect.ValueOf(KeyValue{Key:k, Value:v}))
		}
		fmt.Println("merge    结束")
	}
	return newDataset
}

// 在本地进行排序，形成将使得key相同的处于连续的位置
// 传入的参数为排序的规则函数
func (dataset *DataSet) LocalSort(sortFn interface{}) *DataSet{
	newDataset := NewDataSet(dataset.context, dataset.Type)
	newDataset.SetupShard(len(dataset.Datas))
	step := dataset.context.NewStep()
	step.Name = "LocalSort"
	dataset.context.OneInputForOneOutput(dataset, newDataset, step)

	step.Function = func(task *Task) {
		outChan := task.OutputShards[0].WriteChan
		defer outChan.Close()
		var kvs KvSlice
		var wg sync.WaitGroup
		// 获得所有的kvs
		for _, shard := range task.InputShards{
			//	// 实际上在map这一步的操作中，该task的InputsShard应当只有一个
			//	// 也就是说这个for loop应当执行一次
			wg.Add(1)
			go _handleShaderWriteChannelReceive(&wg, func(v reflect.Value){
				dataset.Lock()
				kvs = append(kvs, v.Interface().(KeyValue))
				dataset.Unlock()
				} , nil,shard)
		}
		wg.Wait()
		fmt.Println("before sorts     ", task.Id ,kvs)
	// 对于获得kvs列表进行排序
		sort.Sort(kvs)
		fmt.Println("after sorts    ", task.Id ,kvs)
		for _,kv := range kvs{
			outChan.Send(reflect.ValueOf(kv))
		}
	}
	return  newDataset
}


func (dataset *DataSet) SetupShard(n int) {
	fmt.Println(dataset)
	ctype := reflect.ChanOf(reflect.BothDir, dataset.Type)
	for i := 0; i < n; i++ {
		ds := &DataShard{
			Id:        i,
			//Parent:    d,
			WriteChan: reflect.MakeChan(ctype, 64), // a buffered chan!
		}
		dataset.Datas = append(dataset.Datas, ds)
	}
}
func (dataset *DataSet)AddOutput(output chan KeyValue){
	fmt.Println(len(dataset.Datas))
	var wg sync.WaitGroup
	defer func() {
		fmt.Println("stop addoutput................")
		close(output)
	}()
	for i, shard := range dataset.Datas{
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, shard *DataShard) {
			defer shard.Close()
			defer wg_.Done()
			for {
				fmt.Println("shard.WriteChan------", i,shard.WriteChan.String())
				if v, OK := shard.WriteChan.Recv(); OK {
					fmt.Println("========+++++++++shard.WriteChan.Recv()　将数据写到channel  output中", v)
					output <- v.Interface().(KeyValue)
				} else {
					fmt.Println("============++++++++++++++++ break")
					break
				}
			}
		}(&wg, shard)
	}
	wg.Wait()
}



func NewDataSet(context *BiuContext, t reflect.Type, ) *DataSet{
	ds := &DataSet{
		Id:len(context.DataSets),
		context:context,
		Type:t,
	}
	context.DataSets = append(context.DataSets, ds)
	return ds
}

// handleV: 处理接受到达数据
// notifyEnd: 通知外部数据接受结束
func _handleShaderWriteChannelReceive(wg *sync.WaitGroup, handleV func(v reflect.Value), notifyEnd func(),shader *DataShard){
	defer wg.Done()
	//defer func() {
	//	shader.Close()
	//}()
	for{
		if v, OK := shader.WriteChan.Recv(); OK{
			//fmt.Println("v:   ", v, "shader ",shader.Id)
			handleV(v)
		}else {
			// 如果有外部通函数，通知外部数据接收结束了
			if notifyEnd != nil{
				notifyEnd()
			}
			break
		}
	}
}