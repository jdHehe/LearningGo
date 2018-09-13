package management

import (
	"sync"
	"net"
)

/*
worker
执行任务
*/
type Worker struct {
	sync.Mutex
	master string // master地址

	listener net.Listener //接受master的任务
}
