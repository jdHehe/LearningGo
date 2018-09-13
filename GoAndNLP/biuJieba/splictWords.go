package biuJieba

import (
	"github.com/adamzy/cedar-go"
	"fmt"
)

//  对字符串进行分割，产生DAG图
func SplitWords(words string, trie *cedar.Cedar) (map[int]int, error){
	words_runes := []rune(words)
	// 映射关系
	words_map := make(map[int]int)
	for i, word_outer := range words_runes{
		fmt.Println("第一个词语", string(word_outer))
		index, _ := trie.Get([]byte(string(word_outer)))
		//if err != nil{
		//	return nil, err
		//}
		if index == 0{
			continue
		}
		words_map[i] = i
		if i == len(words_runes) {
			break
		}
		for j, ex := range words_runes[i+1:]{
			fmt.Println(string(ex))
			index_inner,  _ := trie.Get([]byte(string(words_runes[i:i+2+j])))
			fmt.Println("后一个词语 ", string(words_runes[i:i+2+j]))
			if index_inner == 0{
				fmt.Println(string(words_runes[i:i+2+j]), " 不在整个集合中  ")
				break
			}
			words_map[i] = i+1+j
			fmt.Println("给words_map[",i,"]", "赋值", i+1+j)
		}
	}

	//  从words_map中组合出来所有可能的DAG组合可能
	path := make([]int, 0)
	RecursionSearch(words_map, path, 0, 6)
	return words_map, nil
}

func DynamicProgramming(paths map[int]int, end int){
//	动态规划，寻找最佳的路径
//  不能基于贪心的原则
//  方法一：将所有的可能路径都列举出来，然后选取一条最佳的路径


}
func RecursionSearch(paths map[int]int, path []int,key int, end int) {
	fmt.Println(key, end, path, "  Recusion")
	if key == end{
		path = append(path, end)
		fmt.Println("一条路径", path)
		return
	}
	if key == paths[key] {
		fmt.Println("path = append(path, key)", path, key)
		path = append(path, key)
		key++
		for{
			if key == paths[key]{
				fmt.Println("for path = append(path, key)", path, key)
				path = append(path, key)
				key++
			}else {
				break
			}
		}
		//RecursionSearch(paths, path, key, end)
	}
	for i := key + 1; i <= paths[key]; i++ {
		newPath := DeepCopy(path)
		fmt.Println("新路径  ", i, newPath)
		RecursionSearch(paths, newPath, i, end)
	}

}

func DeepCopy(input[]int)[]int{
	output := make([]int, 0)
	for _, v := range input{
		output = append(output, v)
	}
	return output
}