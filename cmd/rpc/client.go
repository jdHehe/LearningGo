package main

import (

	rpc_local "github.com/jdHeHe/LearningGo/network/rpc"
)


func main(){
	//address := "localhost:1234"
	//serviceMethod := "Arithmetic.Multiply"
	//args := &rpc_local.Args{A:1, B:2}
	//var reply int
	//rpc_local.StartHttpClient(address, serviceMethod, args, &reply)
	//rpc_local.StartTcpClient(address, serviceMethod, args, &reply)
	rpc_local.StartJsonClient()
}


