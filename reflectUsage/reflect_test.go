package reflectUsage

import (
	"testing"
	"reflect"
	"fmt"
)
type helo struct{
	A 	int
	B 	string
}

func (he *helo)SetA(a int){
	he.A = a
}
func (he *helo)AetB(b string){
	he.B = b
}
func (he *helo)GetA()int{
	return he.A
}

/*
he := &helo{1, "12"}
type_ 	:= reflect.TypeOf(*he)
value_ 	:= reflect.ValueOf(&he)  //指针的地址的副本
在reflect包的使用中得注意  值和地址的区别： 例如上面的type和value的取值方式。
值：以值的形式进行反射的时候，获得的是关于这个结构体实例he的内容的一些描述 例如：方法的个数、数据成员的个数
地址：以地址副本的方式使用reflect的时候，获得的是结构体的地址的副本，通过地址可以调用结构体实例he的具体方法call
*/


func TestStruct(t *testing.T){
	he := &helo{1, "12"}
	fmt.Println(he)
	fmt.Println(reflect.TypeOf(he))
	fmt.Println(&he)
	fmt.Println(reflect.TypeOf(&he))
	fmt.Println("========================")
	fmt.Println(reflect.ValueOf(he))
	fmt.Println(reflect.ValueOf(&he))
	type_ 	:= reflect.TypeOf(*he)
	value_object 	:= reflect.ValueOf(*he) // 获得结构体实例he的一个副本的地址
    value_  := reflect.ValueOf(&he)  //获得 结构体实例 he的地址的副本; （传址）这意味这对value_的操作会影响到he的值

	fmt.Println(type_.NumMethod(), type_.NumField())
	fmt.Println(value_object.Kind(), value_object.Type(), value_object.NumField(), value_object.NumMethod())

	//fmt.Println(value_.Field(1))
	heInstanvce := value_.Elem()
	fmt.Println(heInstanvce)
	values := make([]reflect.Value, 1)
	values[0] = reflect.ValueOf(4)
	heInstanvce.MethodByName("SetA").Call(values)

	getA := heInstanvce.MethodByName("GetA")
	A_value := getA.Call(nil)
	fmt.Println(A_value[0])
}
