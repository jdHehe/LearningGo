package biuJieba

import (
	"testing"
	"fmt"
	"os"
)

func TestInitialization(t *testing.T) {
	res, err := Initialization("dict.txt")
	if err != nil{
		fmt.Print(err)
		os.Exit(1)
	}
	//我在北京天安门
	words_list, err := SplitWords("我在东京", res)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("我在北京天安门")
	fmt.Println(words_list)
}

func TestSplitWords(t *testing.T){
	words := "你是谁啊"
	runes := []rune(words)
	fmt.Println(len(runes))
	fmt.Println(string(runes[0]))

	words  = "123344asd"
	runes  = []rune(words)
	fmt.Println(len(runes))
	fmt.Println(string(runes[0]))


}
func TestRecursionSearch(t *testing.T) {

}