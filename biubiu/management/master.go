package management

import (
	"sync"
	"net"
	"net/http"
	"log"
)

/*
初期：
Master 管理worker， 给worker分配任务
*/

type Master struct {
	sync.Mutex
	address string 					// Master的地址
	workers []string				// 注册的workers的地址

	registerChannel chan struct{} 	// worker注册 通知的channel
	doneChannel chan struct{}	  	// 终止当前服务
	l 	net.Listener				// 监听服务注册等

}

/*
运行一个Master		通过http接受各种消息
接受worker的注册
接受客户端的任务
*/
func (master *Master)Run(listenAddressTcp string){
	master_mux := http.NewServeMux()
	master_mux.HandleFunc("/register", master.RegisterHandler)
	master_mux.HandleFunc("/removeWorker", master.RemoveHandler)

	listener, err := net.Listen("tcp", listenAddressTcp)
	if err != nil{
		log.Fatalf("tcp listen address wrong: %v", err)
	}
	err = http.Serve(listener, master_mux)
	if err != nil{
		log.Fatalf("master server cannot server mux:  %v", err)
	}
}

func (master *Master)RegisterHandler(w http.ResponseWriter, r *http.Request){
	// 注册worker
	worker_address := r.Form.Get("address")
	for _,v := range master.workers{
		if v == worker_address{
			return
		}
	}
	master.Mutex.Lock()
	defer master.Mutex.Unlock()
	master.workers = append(master.workers, worker_address)
}
func (master *Master)RemoveHandler(w http.ResponseWriter, r *http.Request){
	worker_address := r.Form.Get("address")
	for i,v := range  master.workers{
		if v== worker_address{
			master.Mutex.Lock()
			master.workers = append(master.workers[:i], master.workers[i+1:]...)
			master.Mutex.Unlock()
			return
		}
	}
}


func NewMaster(address_ string) *Master{
	return &Master{
		address:address_,
		registerChannel:make(chan struct{}),
		doneChannel:make(chan struct{}),
		workers:make([]string, 0),
	}
}
