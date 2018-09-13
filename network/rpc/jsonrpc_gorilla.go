package rpc
import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
	"bytes"
)
/*
client: 将rpc请求的方法以及相关参数序列化成byte数组
		将序列化的数组作为http.Request 的body

		将response的body解码成reply需要的格式

server：注册对Request进行编解码的Codec
		注册Server的Service
		将这个server作为handler注册到路由中
		监听、服务请求，调用server.ServerHttp(w,r)方法
		ServerHttp对请求进行处理，然后将结果编码到response.body中
*/
type Args_Json struct {
	A, B int
}

type Arith_Json int

type Result_Json int

func (t *Arith_Json) Multiply(r *http.Request, args *Args, result *Result_Json) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result_Json(args.A * args.B)
	return nil
}

// 利用gorilla rpc/json
func StartJsonServer(){
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	arith := new(Arith_Json)
	s.RegisterService(arith, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":1234", r)
}

func StartJsonClient(){
	url := "http://localhost:1234/rpc"
	args := &Args_Json{
		A: 2,
		B: 3,
	}
	// 构造一个jsonRpc的request 并且通过json.Marshal()对数据进行编码
	message, err := json.EncodeClientRequest("Arith_Json.Multiply", args)
	if err != nil {
		log.Fatalf("%s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		log.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error in sending request to %s. %s", url, err)
	}
	defer resp.Body.Close()

	var result Result_Json
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		log.Fatalf("Couldn't decode response. %s", err)
	}
	log.Printf("%d*%d=%d\n", args.A, args.B, result)
}