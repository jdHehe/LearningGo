package rpc

import (
	"reflect"
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"math/rand"
	"time"
	"strings"
	"fmt"
)

type reqMsg struct {
	endName 	interface{} 	//终端名称
	svcMeth 	string			//请求的服务的方法名称
	argsType	reflect.Type	//请求的参数类型
	args 		[]byte			//请求的参数的值
	replyCh 	chan replyMsg	//回复的通道
}
type replyMsg struct {
	ok bool
	reply []byte
}
type ClientEnd struct {
	endName	interface{}
	ch 		chan reqMsg
}
func (e *ClientEnd)Call(svcMeth string, args interface{}, reply interface{}) bool{
	// 构造rpc请求的消息
	req := reqMsg{}
	req.endName = e.endName
	req.svcMeth = svcMeth
	req.argsType = reflect.TypeOf(args)
	req.replyCh = make(chan replyMsg)

	// 对参数args数据进行编码
	qb := new(bytes.Buffer)
	qe := gob.NewEncoder(qb)
	qe.Encode(args)
	req.args = qb.Bytes()
	e.ch <- req

	// 等待回复
	rep := <-req.replyCh
	if rep.ok {
		rb := bytes.NewBuffer(rep.reply)
		rd := gob.NewDecoder(rb)
		if err := rd.Decode(reply); err != nil{
			log.Fatalf("ClientEnd.call(): decode error")
		}
		return true
	}else {
		return false
	}
}

type Network struct {
	mu 				sync.Mutex
	reliable		bool
	longDelays		bool
	longRecordering	bool
	ends 			map[interface{}]*ClientEnd
	enabled			map[interface{}]bool
	servers			map[interface{}]*Server
	connections  	map[interface{}]interface{}
	endch			chan reqMsg
}

func MakeNetwork() *Network{
	// 创建一个Network对象，并进行相关配置
	rn := &Network{}
	rn.reliable = true
	rn.ends = map[interface{}]*ClientEnd{}
	rn.enabled = map[interface{}]bool{}
	rn.servers = map[interface{}]*Server{}
	rn.connections = map[interface{}](interface{}){}
	rn.endch = make(chan reqMsg)

	//	通过channel接受rpc请求的信息，对每个请求创建一个routine
	go func() {
		for xreq := range rn.endch{
			fmt.Println("get msg", xreq)
			go rn.ProcessReq(xreq)
		}
	}()
	return rn
}

// 处理rpc 请求
// to be continue
func (network *Network)ProcessReq(req reqMsg) {
	enabled, servername, server, reliable, longreordering := network.ReadEndnameInfo(req.endName)
	if enabled && servername != nil && server != nil {
		if reliable == false {
			ms := rand.Int() % 27
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}
		if reliable == false && (rand.Int()%1000) < 100 {
			req.replyCh <- replyMsg{false, nil}
			return
		}

		ech := make(chan replyMsg)
		go func() {
			// 处理请求，获得返回结果
			r := server.dispatch(req)
			fmt.Println("after dispatch", r)
			ech <- r
		}()

		var reply replyMsg
		replyOK := false
		serverDead := false
		for replyOK == false && serverDead == false {
			select {
			case reply = <-ech:
				fmt.Println("set replyOK true")
				replyOK = true
			case <-time.After(100 * time.Microsecond):
				serverDead = network.IsServerDead(req.endName, servername, server)
			}
		}
		fmt.Println("========================")
		fmt.Println(serverDead, replyOK)
		serverDead = network.IsServerDead(req.endName, servername, server)
		fmt.Println("========================")
		fmt.Println(serverDead, replyOK)
		if replyOK == false || serverDead == true{
			req.replyCh <- replyMsg{false, nil}

		}else if reliable == false && (rand.Int()%1000) < 100{
			req.replyCh <- replyMsg{false, nil}
		}else if longreordering==true && rand.Intn(900)<600 {
			ms := 200 + rand.Intn(1+rand.Intn(2000))
			time.Sleep(time.Duration(ms) * time.Millisecond)
			req.replyCh <- reply
		}else {
			req.replyCh <- reply
		}
	}else {
		ms := 0
		if network.longDelays{
			ms = (rand.Int() % 7000)
		}else {
			ms = (rand.Int()%100)
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		req.replyCh <- replyMsg{false, nil}
	}
}
func (network *Network)IsServerDead(endName interface{}, serverName interface{}, server *Server) bool{
	network.mu.Lock()
	defer network.mu.Unlock()

	if network.enabled[endName] == false || network.servers[serverName] != server{
		return  true
	}
	return false
}

// 将请求消息reqMsg 元素提取出来
func (network *Network) ReadEndnameInfo(endname interface{}) (enabled bool,
	servername interface{}, server *Server, reliable bool, longreordering bool,
) {
	network.mu.Lock()
	defer network.mu.Unlock()

	enabled = network.enabled[endname]
	servername = network.connections[endname]
	if servername != nil{
		server = network.servers[servername]
	}
	reliable = network.reliable
	longreordering = network.longRecordering
	return
}

// 创建一个终端
func (network *Network)MakeEnd(endName interface{}) *ClientEnd{
	network.mu.Lock()
	defer network.mu.Unlock()
	if _, ok := network.ends[endName]; ok{
		log.Fatalf("MakeEnd: %v already exists\n", endName)
	}

	e := &ClientEnd{}
	e.endName = endName
	e.ch = network.endch
	network.ends[endName] = e
	network.enabled[endName] = false
	network.connections[endName] = nil

	return e
}

// 向network 添加/删除 一个server
func (network *Network) AddServer(serverName interface{}, rs *Server){
	network.mu.Lock()
	defer  network.mu.Unlock()
	network.servers[serverName] = rs
}
func (network *Network) DeleteServer(serverName interface{}){
	network.mu.Lock()
	defer  network.mu.Unlock()
	network.servers[serverName] = nil
}
// 将终端和server连接起来，并记录这个连接
func (network *Network) Connect(endName interface{}, serverName interface{}){
	network.mu.Lock()
	defer network.mu.Unlock()
	network.connections[endName] = serverName
}
// 控制终端是否可获得
func (network *Network)Enable(endName interface{}, enabled bool){
	network.mu.Lock()
	defer network.mu.Unlock()

	network.enabled[endName] = enabled
}

// server 可以进行多种service
// 通过dispatch分发不同的rpc请求任务
type Server struct {
	mu 			sync.Mutex
	services	map[string]*Service  // 服务名和服务的映射
	count 		int
}
// 新建Server， Server包括多个对应不同服务的Service
func MakeServer() *Server{
	rs := &Server{}
	rs.services = map[string]*Service{}
	return rs
}
func (server *Server)AddService(service *Service){
	server.mu.Lock()
	defer server.mu.Unlock()
	server.services[service.name] = service
}

//	进行请求的分发
func (server *Server) dispatch(req reqMsg) replyMsg {
	server.mu.Lock()
	server.count += 1
	// 获得method和service的相关信息
	dot := strings.LastIndex(req.svcMeth, ".")
	serviceName := req.svcMeth[:dot]
	methodName := req.svcMeth[dot+1:]
	service, ok := server.services[serviceName]
	server.mu.Unlock()

	if ok{
		return service.dispatch(methodName, req)
	}else{
		choices := []string{}
		for k, _ := range server.services{
			choices = append(choices, k)
		}
		return replyMsg{false, nil}
	}

}

// Service 是rpc具体调用的方法的载体
// 最终server将请求的执行分发到具体的service
type Service struct {
	name 	string
	rcvr	reflect.Value
	typ 	reflect.Type
	methods map[string]reflect.Method
}

func MakeService(rcvr interface{}) *Service{
	svc := &Service{}
	svc.typ = reflect.TypeOf(rcvr)
	svc.rcvr = reflect.ValueOf(rcvr)
	svc.name = reflect.Indirect(svc.rcvr).Type().Name()
	svc.methods = map[string]reflect.Method{}

	for  m := 0; m <svc.typ.NumMethod(); m++{
		method 	:= svc.typ.Method(m)
		mtype 	:= method.Type
		mname 	:= method.Name
		if method.PkgPath != "" || mtype.NumIn() != 3 || mtype.In(2).Kind()!= reflect.Ptr||
			mtype.NumOut() != 0{

		}else {
			svc.methods[mname] = method
		}
	}
	return svc
}

//  调用特定的service(具体方法的载体) 的方法
func (service *Service) dispatch(methodName string, req reqMsg)  replyMsg {
	if method, ok := service.methods[methodName]; ok{
		// 解码 参数args
		args := reflect.New(req.argsType)
		ab := bytes.NewBuffer(req.args)
		ad := gob.NewDecoder(ab)
		ad.Decode(args.Interface())

		// 通过反射获得注册的Service的方法的返回值reply的类型，通过类型构造出新得reply的值，用于call方法
		replyType := method.Type.In(2)
		replyType = replyType.Elem()
		replyValue := reflect.New(replyType)

		function := method.Func
		function.Call([]reflect.Value{service.rcvr, args.Elem(), replyValue})

		rb := new(bytes.Buffer)
		re := gob.NewEncoder(rb)
		re.EncodeValue(replyValue)

		return replyMsg{true, rb.Bytes()}
	}else {
		choices := []string{}
		for k, _ := range service.methods{
			choices =append(choices, k)
		}
		return replyMsg{false, nil}
	}
}