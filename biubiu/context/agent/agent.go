package agent

import (
	"net/http"
	"strconv"
	"net"
)

type Agent struct {

}


//  通过http请求接受 task任务


//  处理具体的任务
func (agent *Agent) DoStageJob(path string, stageGroupId int, stageId int){

}
func (agent *Agent) NewAgentServer(address string) (*AgentServer, error){
	listener, err := net.Listen("tcp", address)
	if err != nil{
		return nil, err
	}
	return &AgentServer{
		agent:agent,
		listener:listener,
	}, nil
}