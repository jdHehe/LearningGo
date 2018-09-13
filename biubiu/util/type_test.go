package util

import (
	"testing"
	"sync"
	"fmt"
)

func TestMap(t *testing.T){
	var wg sync.WaitGroup
	var lock sync.Mutex
	m := make(map[int]int)
	for i:= 0; i<2; i++{
		wg.Add(1)
		go func(mp map[int]int) {
			defer wg.Done()
			for j:= 0; j<5; j++{
				lock.Lock()
				if v,OK := mp[j]; OK{
					mp[j]= v+j
				}else {
					mp[j] = j
				}
				lock.Unlock()
			}
		}(m)
	}
	wg.Wait()

	fmt.Println("map ",m)
}