package context

import (
	"reflect"
	"fmt"
)

// 数据分片
type DataShard struct {
	FilePath string
	WriteChan reflect.Value
	Id 	int
	AgentLocaltion string //拥有这个shader的agent的位置
}

func (shader *DataShard)Close(){
	fmt.Println("close shard ", shader.Id)
	shader.WriteChan.Close()
}
