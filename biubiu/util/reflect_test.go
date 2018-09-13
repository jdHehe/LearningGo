package util

import (
	"testing"
	"fmt"
)

func TestFuncType(t *testing.T) {
	type_ := FuncType(func(chan string) {})
	fmt.Println(type_)
}
