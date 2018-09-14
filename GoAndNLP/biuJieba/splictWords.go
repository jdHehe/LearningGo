// create by kakunka

package biuJieba

import (
	"github.com/adamzy/cedar-go" // double array trie 高效的搜索的树结构
	"fmt"
	"container/list"
)

type WordFreqToIndex struct {
	Freq 	int 	// 词频
	index 	int		// 这个词的结尾的单词处在词组中的位置
}

//  对字符串进行分割，产生DAG图
func SplitWords(words string, trie *cedar.Cedar) (map[int]int, error){
	words_runes := []rune(words)
	// 映射关系
	words_map := make(map[int]int)
	for i, word_outer := range words_runes{
		index, _ := trie.Get([]byte(string(word_outer)))
		fmt.Println("index")
		fmt.Println(index)
		if index == 0{
			continue
		}
		words_map[i] = i
		if i == len(words_runes) {
			break
		}
		for j, _ := range words_runes[i+1:]{
			index_inner,  _ := trie.Get([]byte(string(words_runes[i:i+2+j])))
			if index_inner == 0{
				break
			}
			words_map[i] = i+1+j
		}
	}

	fmt.Println("word_map")
	fmt.Println(words_map)
	route := DynamicProgramming(words_map, words, trie)
	fmt.Println("route...................")
	fmt.Println(route)
	res := new(list.List)
	index := 0
	for index < len(words){
		end := route[index].index
		word := string(words_runes[index:end+1])
		res.PushBack(word)
		index = end + 1
	}
	fmt.Println("res..................")
	fmt.Println(res.Len())
	for e := res.Front(); e != nil; e = e.Next(){
		fmt.Print(e.Value, " ")
	}
	fmt.Println()
	return words_map, nil
}

func DynamicProgramming(DAG map[int]int,  words string, trie *cedar.Cedar) map[int]WordFreqToIndex{
//	动态规划，寻找最佳的路径
//  不能基于贪心的原则
//  方法一：将所有的可能路径都列举出来，然后选取一条最佳的路径
//  最大概率路径问题

//  从后往前遍历words，寻找概率最大的组合
	N := len([]rune(words))
	fmt.Println("参数")
    fmt.Println(DAG, N, words)
	words_runes := []rune(words)
	route := make(map[int]WordFreqToIndex)
	route[N] = WordFreqToIndex{Freq:0, index:N}
	var end  =  0
	for i:=N-1; i>=0; i-- {
		max := 0
		end = i
		for j:=i; j<= DAG[i]; j++{
			value, _ := trie.Get([]byte(string(words_runes[i:j+1])))
			fmt.Println("词语  ", string(words_runes[i:j+1]), value, i, j)
			value += route[j+1].Freq
			if max < value {
				max = value
				end = j
			}
		}
		route[i] =  WordFreqToIndex{
			Freq:  max,
			index: end,
		}
	}

	return route
}

// 尝试从前往后穷尽所有可能的路径，寻找概率最大的那条路径
//func RecursionSearch(paths map[int]int, path []int,key int, end int) {
//	if key == end{
//		path = append(path, end)
//		return
//	}
//	if key == paths[key] {
//		path = append(path, key)
//		key++
//		for{
//			if key == paths[key]{
//				path = append(path, key)
//				key++
//			}else {
//				break
//			}
//		}
//	}
//	for i := key + 1; i <= paths[key]; i++ {
//		newPath := DeepCopy(path)
//		RecursionSearch(paths, newPath, i, end)
//	}
//
//}

func DeepCopy(input[]int)[]int{
	output := make([]int, 0)
	for _, v := range input{
		output = append(output, v)
	}
	return output
}