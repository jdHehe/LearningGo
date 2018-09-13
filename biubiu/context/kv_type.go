package context

import "reflect"

type KeyValue struct {
	Key  	interface{}  //{String() string}
	Value	interface{}

}


var (
	KeyValueType = reflect.TypeOf(KeyValue{})
)