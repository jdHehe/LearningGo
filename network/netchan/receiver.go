package netchan

import "net"

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
