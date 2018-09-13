package netchan

import (
	"net"
	//"reflect"
	"encoding/gob"
	"github.com/jdHeHe/LearningGo/biubiu/context"
	"fmt"
)
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

func GetSenderChannelByReflectValue(target string, Instance interface{}) (chan interface{}, error){
	_, err := net.ResolveTCPAddr("tcp", target)
	if err != nil{
		return nil, err
	}
	conn, err := net.Dial("tcp", target)
	enc := gob.NewEncoder(conn)
	if err != nil{
		return nil, err
	}
//	msgs := make(chan reflect.Value)
	msgs := make(chan interface{})
	go func() {
		for msg := range msgs{
			switch Instance.(type) {
			case context.KeyValue:
				value := msg.(context.KeyValue)
				err = enc.Encode(&value)
				if err != nil{
					fmt.Println(err)
					break
				}
			case string:
				value := msg.(string)
				err = enc.Encode(&value)
				if err != nil{
					fmt.Println(err)
					break
				}
			}
		}
	}()
	return msgs, err
}