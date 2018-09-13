package netchan

import "net"
/*
Get a channel to send  message to target server
*/
func GetSenderChannel(target string) (chan Message, error){
	_, err := net.ResolveTCPAddr("tcp", target)
	if err != nil{
		return nil, err
	}
	conn, err := net.Dial("tcp", target)
	if err != nil{
		return nil, err
	}
	msgs := make(chan Message)

	for msg := range msgs{
		conn.Write(msg.Bytes())
	}
	return msgs, nil
}