package main

import (
	"fmt"
	"runtime/debug"
)

//	定义一种错误类型，包含尽量多的信息
//  通过wrapError 来包装成MyError错误类型，
type MyError struct {
	Inner 		error
	Message 	string
	StackTrace 	string
	Misc 		map[string]interface{}
}
func wrapError(err error, messagef string, msgArgs ... interface{}) MyError{
	return MyError{
		Inner:		err,
		Message:	fmt.Sprintf(messagef, msgArgs...),
		StackTrace:	string(debug.Stack()),
		Misc:		make(map[string]interface{}),
	}
}

// Timeout & Cancellation
// Timeout  :系统性能饱和，无法处理请求； 对过期的资源进行处理
// Cancellation :取消的可能有 超时、用户干预、父进程的主动取消、重复的请求导致的取消
