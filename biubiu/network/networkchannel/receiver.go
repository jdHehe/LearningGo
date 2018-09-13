package netchan

import (
	"net"
	context "github.com/jdHeHe/LearningGo/biubiu/context"
	"encoding/gob"
	"fmt"
)

// 返回接受 网络消息的channel
func GetReceiverChannel(target string)(chan Message, error){
	conn, err := net.Dial("tcp", target)
	if err != nil{
		return nil, err
	}
	defer conn.Close()
	recevie := make(chan Message)
	buff := make([]byte, 0)
	for{
		_, err := conn.Read(buff)
		if err != nil {
			break
		}
		msg := new(Message)
		msg.ToMsg(buff)
		recevie <- *msg
	}
	return recevie, nil
}

func GetReceiverChannelByReflectValue(target string, Instance interface{})(chan interface{}, error){
	listener, err := net.Listen("tcp", target)

	if err != nil{
		return nil, err
	}
	recevie := make(chan interface{})
	//buff := make([]byte, 0)
	go func() {
		defer listener.Close()
		for {
			fmt.Println("开始接受连接")
			conn, err := listener.Accept()
			if err != nil{
				fmt.Println("err: || listener.Accept()", err)
			}
			go HandleConn(Instance, conn, recevie)
		}
	}()
	fmt.Println("get out of GetReceiverChannelByReflectValue ")
	return recevie, nil
}
func HandleConn(Instance interface{}, conn net.Conn, receive chan interface{}){
	fmt.Println("接收到 连接")
	dec := gob.NewDecoder(conn)
	for {
		switch Instance.(type) {
		case context.KeyValue:
			value := context.KeyValue{}
			err := dec.Decode(&value)
			if err != nil {
				fmt.Println(err)
				break
			}
			receive <- value
		case string:
			value := ""
			fmt.Println("dec:   ", dec)
			err := dec.Decode(&value)
			if err != nil {
				fmt.Println(err)
				break
			}
			receive <- value
		default:
			fmt.Println("接收到的类型不是预定的类型，无法处理")
		}
	}
}