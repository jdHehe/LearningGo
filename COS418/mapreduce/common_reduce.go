package mapreduce

import (
	"encoding/json"
	"os"
	"fmt"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()

	// 每个Reducer会接受多个Mapper产生的文件
	//fmt.Println("doReduce  common_reduce.go. ReduceTasknNumber is", reduceTaskNumber, " nMap is", nMap)
	var kv KeyValue
	kvs := make(map[string] []string)
	for i:=0; i<nMap; i++{
		fileName := reduceName(jobName, i, reduceTaskNumber)
		//fmt.Println("fileName of reduceName produced by map ", fileName)
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil{
			panic(err)
		}
		decoder := json.NewDecoder(file)
		for {
			err := decoder.Decode(&kv)
			if err != nil{
				fmt.Println(err)
				break
			}
			kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
		}
	}
	//fmt.Println(kvs)
	mergerFile := mergeName(jobName, reduceTaskNumber)
	fmt.Println("mergerFile is ",mergerFile)
	file_merger, err :=  os.Create(mergerFile)
	defer file_merger.Close()
	if err != nil{
		panic(err)
	}
	enc := json.NewEncoder(file_merger)
	for k, kv := range kvs{
		enc.Encode(KeyValue{k, reduceF(k, kv)})
	}
}
