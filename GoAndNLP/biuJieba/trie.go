package biuJieba

import (
	"github.com/adamzy/cedar-go"
	"os"
	"bufio"
	"io"
	"bytes"
	"strings"
	"strconv"
)


func Initialization(fileName string) (*cedar.Cedar, error){
	trie := cedar.New()
	file , err := os.Open(fileName)
	if err != nil{
		return nil, err
	}
	br := bufio.NewReader(file)
	words_frequence := 0
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF{
			break
		}
		line_string := bytes.NewBuffer(line).String()
		string_slices := strings.Split(line_string, " ")
		if len(string_slices) != 3{
			continue
		}
		words_frequence, err  = strconv.Atoi(string_slices[1])
		if err != nil{
			return nil, err
		}
		trie.Insert(bytes.NewBufferString(string_slices[0]).Bytes(), words_frequence)
	}
	return trie, nil
}
