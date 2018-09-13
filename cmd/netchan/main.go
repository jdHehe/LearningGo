package main

import (
	//"fmt"
	//"flag"
	//"net/http"
	//"bytes"
	//"github.com/jdHeHe/LearningGo/network/netchan"
	"encoding/json"
	"os"
	"fmt"
)

func main(){
//	isServer := flag.Bool("server", false, "start as a server")
//	address  := flag.String("address", "localhost:3000", "server address and port")
//	flag.Parse()
//	if *isServer {
//		fmt.Println(*address)
//		http.HandleFunc("/topic", handleClientRequest)
//		http.ListenAndServe(*address, nil)
//	}else {
//		fmt.Println("not server")
//	}

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			fmt.Println(err)
			return
		}
		for k := range v {
			if k != "Name" {
				delete(v, k)
			}
		}
		if err := enc.Encode(&v); err != nil {
			fmt.Println(err)
		}
	}
}
//
//type clientHandler struct {
//
//}
//func handleClientRequest(w http.ResponseWriter, r *http.Request){
//	address := r.URL.Query()["address"]
//	fmt.Println(address)
//	msg, err := netchan.GetReceiverChannel(address[0])
//	w.Write(bytes.NewBufferString(address[0]).Bytes())
//	make()
//}