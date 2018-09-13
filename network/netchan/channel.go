package netchan

import "bytes"

type Message struct {
	content string
}
func (msg * Message)Bytes() []byte{
	return  bytes.NewBufferString(msg.content).Bytes()
}
func (msg *Message)ToMsg(message []byte){
	msg.content = bytes.NewBuffer(message).String()
}
