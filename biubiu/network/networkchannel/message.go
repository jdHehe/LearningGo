package netchan

type Message struct {
	content []byte
}
func (msg * Message)Bytes() []byte{
	return  msg.content
}
func (msg *Message)ToMsg(message []byte){
	msg.content = message
}
