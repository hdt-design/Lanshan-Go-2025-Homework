package main

import (
	"fmt"
	"sync"
)

func main() {
	sum := 0
	lock := sync.Mutex{}
	for range 10 {
		go func() {
			for range 100000 {
				lock.Lock()
				sum += 1
				lock.Unlock()
			}
		}()
	}
	fmt.Println(sum)
}
