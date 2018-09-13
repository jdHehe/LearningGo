package main

import (
	"fmt"
	"bytes"
	//"reflect"
	//"math"
)

//  一致性hash
//	DHT
//  主要是计算字符串之间的距离，然后根据距离为资源选择合适的Node

type Node struct {
	address []byte
}

type Resource struct {
	content string
	address []byte
	NodeId int
}

func main(){
	res1 := Resource{
		address:bytes.NewBufferString("1231231234").Bytes(),
	}
	res2 := Resource{
		address:bytes.NewBufferString("4234231234").Bytes(),
	}
	node1 := Node{
		address:bytes.NewBufferString("1111231234").Bytes(),
	}
	node2 := Node{
		address:bytes.NewBufferString("1114231234").Bytes(),
	}
	ns := []Node{node1, node2}
	id, err := FindNode(res1, ns)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("res1, id ",id)

	id , err = FindNode(res2, ns)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("res2, id ",id)
}

func FindNode(r Resource,ns []Node) (int,error){
	nearNode := make([]int, 0)
	Interval := make([][]uint8, len(ns))

	// 算出资源到所有节点的距离
	for j, n := range ns{
		interval := make([]uint8, len(n.address))
		for i, a := range n.address{
			interval[i] = MinusUint8(r.address[i], a)
		}
		Interval[j] = interval
	}
	// 选择距离最短的
	// 从前往后  每一个byte的比较
	var smallLestByte uint8 = 255
	for  i := 0 ;i<len(r.address); i++{
		smallLestByte = 255
		//   只需要遍历前置位相同的名单中的节点
		if (len(nearNode) > 0){
			for _, nodeIndex := range nearNode{
				if (smallLestByte > Interval[nodeIndex][i]){
					smallLestByte = Interval[nodeIndex][i]
					nearNode = make([]int, 0)
					nearNode = append(nearNode, nodeIndex)
				}else if (smallLestByte == Interval[nodeIndex][i]){
					nearNode = append(nearNode, nodeIndex)
				}
			}
		}else {
			//  只有第一轮需要遍历所有的距离来找到较小的距离
			for nodeIndex, distance := range Interval {
				if (smallLestByte > distance[i]) {
					smallLestByte = distance[i]
					nearNode = make([]int, 0)
					nearNode = append(nearNode, nodeIndex)
				} else if (smallLestByte == distance[i]) {
					nearNode = append(nearNode, nodeIndex)
				}
			}
		}
		if (len(nearNode) == 1){
			return nearNode[0],nil
		}
	}
	return  -1 , nil
}

func MinusUint8( a uint8,b uint8)uint8{
	if a>b{
		return a - b
	}else{
		return b - a
	}
}
