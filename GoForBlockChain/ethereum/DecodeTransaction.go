package main

import (
	rlp "github.com/ethereum/go-ethereum/rlp"
	"strings"
	"fmt"
)

func main(){
	fmt.Println("==========")
	con := DecodeTx("0x6449ca9f2764e4334b88cb571a02c3517e2f0fe02545eb5767bd26cbd67c74cb")
	fmt.Print(con)

	fmt.Println(DecodeTx("0x6449ca9f2764e4334b88cb571a02c3517e2f0fe02545eb5767bd26cbd67c74cb"))






}


func DecodeTx(content string) string{
	reader := strings.NewReader(content)
	var result string
	rlp.Decode(reader, &result)
	return  result
}
