package util

import (
	"reflect"
)

/*
反射： 做类型判断
*/

//  获取函数fn的输出类型
// 	如果channel的最后一个参数为channel，说明这个函数是通过channel将返回值传出的
//
func FuncType(fn interface{}) reflect.Type{
	funcEle := reflect.TypeOf(fn)
	// 入参为chan类型，返回这个类型
	//fmt.Println(funcEle.In(funcEle.NumIn()-1).Elem(), "hhh")
	if funcEle.In(funcEle.NumIn()-1).Kind() == reflect.Chan{
		//fmt.Println(funcEle.In(funcEle.NumIn()-1).Elem())
		return funcEle.In(funcEle.NumIn()-1).Elem()
	}else if funcEle.NumOut()==1{
		// 如果返回值为一个，则返回这个值
		//fmt.Println( funcEle.Out(funcEle.NumOut()-1))
		return funcEle.Out(funcEle.NumOut()-1)
	}
	return nil
}
