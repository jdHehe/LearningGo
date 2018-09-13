package main

import (
	"github.com/jdHeHe/LearningGo/network/rpc"
	"fmt"
	"time"
	"sync"
	"strconv"
)

type JunkArgs struct {
	X int
}
type JunkReply struct {
	X string
}

type JunkServer struct {
	mu   sync.Mutex
	log1 []string
	log2 []int
}

func (js *JunkServer) Handler1(args string, reply *int) {
	js.mu.Lock()
	defer js.mu.Unlock()
	js.log1 = append(js.log1, args)
	*reply, _ = strconv.Atoi(args)
}

func (js *JunkServer) Handler2(args int, reply *string) {
	js.mu.Lock()
	defer js.mu.Unlock()
	js.log2 = append(js.log2, args)
	*reply = "handler2-" + strconv.Itoa(args)
}

func (js *JunkServer) Handler3(args int, reply *int) {
	js.mu.Lock()
	defer js.mu.Unlock()
	time.Sleep(20 * time.Second)
	*reply = -args
}

// args is a pointer
func (js *JunkServer) Handler4(args *JunkArgs, reply *JunkReply) {
	reply.X = "pointer"
}

// args is a not pointer
func (js *JunkServer) Handler5(args JunkArgs, reply *JunkReply) {
	reply.X = "no pointer"
}



type Student struct {
	name string
	age int
}
type StudentArgs struct{
	name string
	age int
}


func (student *Student)SetName(args *StudentArgs, reply*string){
	fmt.Println(args)
	student.name = args.name
	*reply = "setname"
}
func (student *Student)SetAge(args *StudentArgs, reply *int){
	student.age = args.age
	*reply = 2
}
func (student *Student)GetName(args *StudentArgs, reply *string) {
	*reply = student.name
}
func (student *Student)GetAge(args *StudentArgs ,reply *int){
 	*reply = student.age
}


func main() {
	network := rpc.MakeNetwork()
	highschool := network.MakeEnd("highschool")

	// 创建service
	student := &Student{name:"ligang", age:12}
	student_service := rpc.MakeService(student)

	server := rpc.MakeServer()
	server.AddService(student_service)

	network.AddServer("student", server)
	network.Connect("highschool", "student")
	network.Enable("highschool", true)
	{
		reply := ""
		args := StudentArgs{name: "lee", age: 233}
		highschool.Call("Student.SetName", &args, &reply)
		fmt.Println(reply)
	}
	{
		reply_inner := ""
		highschool.Call("Student.GetName",nil, &reply_inner)
		fmt.Println(reply_inner)
	}
	{
		reply_age := 0
		highschool.Call("Student.GetAge", nil, &reply_age)
		fmt.Println(reply_age)
	}
}
