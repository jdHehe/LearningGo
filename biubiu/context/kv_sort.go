package context

import (
	"strings"
)

type KeyValueTest struct {
	Key string
	Value int
}

type KvSlice []KeyValue

func (kvs KvSlice)Len()int{
	return len(kvs)
}

func (kvs KvSlice)Swap(i, j int){
	kvs[i], kvs[j] = kvs[j], kvs[i]
}
func (kvs KvSlice)Less(i, j int) bool{
	res := strings.Compare(kvs[i].Key.(string), kvs[j].Key.(string))
	return res >= 0
}
