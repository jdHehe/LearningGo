package agent

import (
	"net/http"
	"strconv"
	"net"
)

type AgentServer struct {
	agent *Agent
	listener net.Listener
}

func (server *AgentServer)HandleTaskReuqest(w  http.ResponseWriter, r* http.Request){
	path :=	r.PostForm.Get("path")
	stageGroupIdStr := r.PostForm.Get("stageGroupId")
	stageIdStr 	 := r.PostForm.Get("stageId")
	stageGroupId, err := strconv.Atoi(stageGroupIdStr)
	stageId, err 	 := strconv.Atoi(stageIdStr)
	if err != nil{
		w.WriteHeader(400)
	}
	server.agent.DoStageJob(path, stageGroupId, stageId)
}

func (server *AgentServer)Run(){
	http.serve
}