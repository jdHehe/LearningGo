package main

import rpc_local "github.com/jdHeHe/LearningGo/network/rpc"

func main()  {
	go rpc_local.StartServer()
	rpc_local.CallRpcBySynchronous()
	rpc_local.CallRpcByAsynchronous()
}
