package main

import (
	local_rpc "github.com/jdHeHe/LearningGo/network/rpc"
	"fmt"
	"net/rpc"
	"net"
	"net/http"
)

func main(){
	//rpc_local.StartHttpServer("localhost:1234")
	//rpc_local.StartTcpServer("localhost:1234")
	local_rpc.StartJsonServer()
}
func ServerHttp(){
	arith := new(local_rpc.Arith)
	server := rpc.NewServer()
	server.RegisterName("Arithmetic", arith)
	server.HandleHTTP("/", "debug")
	l, e := net.Listen("tcp", ":1234")
	if e != nil{
		fmt.Println("listen error:", e)
	}
	http.Serve(l ,nil)
}
func ServerTcp(){
	arith := new(local_rpc.Arith)
	server := rpc.NewServer()
	server.RegisterName("Arithmetic", arith)
	l, e := net.Listen("tcp", ":1234")
	if e != nil{
		fmt.Println("listen error:", e)
	}
	server.Accept(l)
}