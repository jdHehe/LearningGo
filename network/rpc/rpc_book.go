package rpc

import (
	"github.com/kataras/iris/core/errors"
	"net/rpc"
	"net"
	"net/http"
	"fmt"
)

// service
type Args struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int

func (t *Arith)Multiply(args *Args, reply *int) error{
	*reply = args.A * args.B
	return  nil
}
func (t *Arith)Divide(args *Args, quo*Quotient) error{
	if args.B == 0{
		return errors.New("divided by zero")
	}
	quo.Quo = args.A/args.B
	quo.Rem = args.A%args.B
	return nil
}

func StartHttpServer(address string){
	arith := new(Arith)
	server := rpc.NewServer()
	// 注册服务
	server.RegisterName("Arithmetic", arith)
	server.HandleHTTP("/", "debug")
	l, err := net.Listen("tcp", address)
	if err != nil{
		panic("listen error")
	}
	http.Serve(l ,nil)
}
func  StartTcpServer(address string){
	arith := new(Arith)
	server := rpc.NewServer()
	server.RegisterName("Arithmetic", arith)
	l, e := net.Listen("tcp", ":1234")
	if e != nil{
		fmt.Println("listen error:", e)
	}
	server.Accept(l)
}
// rpc客户端
func StartTcpClient(address string, serviceMethod string, args *Args, reply * int){
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil{
		fmt.Println("dialing:", err)
	}
	client := rpc.NewClient(conn)
	err  = client.Call("Arithmetic.Multiply", args, &reply)
	if err != nil{
		fmt.Println("dialing:", err)
	}
	fmt.Println(reply)
}
func StartHttpClient(address string, serviceMethod string, args *Args, reply * int){
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil{
		fmt.Println("dialing:", err)
	}
	err = client.Call(serviceMethod, args, &reply)
	if err != nil{
		panic(err)
	}
	fmt.Println(reply)
}