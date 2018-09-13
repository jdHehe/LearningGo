package main

import (
	ethrpc "github.com/onrik/ethrpc"
	"sync"
	"log"
	"fmt"
)

var once sync.Once
var EtherRpcClient *ethrpc.EthRPC

func init() {
	once.Do(func() {
		EtherRpcClient = ethrpc.New( "http://192.168.1.114:8501" )
	})
}

func GetTransactionByHash(hash string) *ethrpc.Transaction{
	 tx, err := EtherRpcClient.EthGetTransactionByHash(hash)
	 if err != nil{
	 	log.Fatal(err)
	 }
	 return tx
}

func GetBlockTransactionsByNumber(number int) []ethrpc.Transaction{
	block, err := EtherRpcClient.EthGetBlockByNumber(number ,true)
	if err != nil{
		log.Fatal(err)
	}
	return block.Transactions
}

func GetBlockNumber() int{
	number, err := EtherRpcClient.EthBlockNumber()
	if err != nil{
		log.Fatal(err)
	}
	return number
}

const interationNumber   = 100
func main(){
	//hash := "0x1309ddbb4f5f3a030985eb05f3ced1f781bfae69fffbc237b785b9afc4a4c6c6"
	//tx := GetTransactionByHash(hash)
	//fmt.Printf("from: %s\n to: %s\n amount: %s",tx.From, tx.To, tx.Value)
	var wg sync.WaitGroup

	blockNumber := GetBlockNumber()
	routines := blockNumber / interationNumber
	wg.Add(routines)
	for i:=0; i< routines; i++{
		fmt.Printf("第 %d 个轮次 \n",i)
		go func() {
		// 每个goroutine查询一段块
		defer wg.Done()
		for  j:=i*100; j<(i+1)*100; j++{
			for _,tx := range GetBlockTransactionsByNumber(j){
				fmt.Printf("from: %s  to: %s  amount: %b \n",tx.From, tx.To, tx.Value)
			}
			fmt.Println()
		}
		}()
	}
	wg.Wait()

	//txs := GetBlockTransactionsByNumber(10)
	//for _,v := range txs{
	//	fmt.Print(v)
	//	fmt.Println("=====================")
	//}

}