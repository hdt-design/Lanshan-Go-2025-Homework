package main

import (
	"fmt"
	"sync"
)

type task struct {
	runnable func()
}

func main() {
	sum := 0
	lock := sync.Mutex{}
	taskch := make(chan task, 100)

	for i := 0; i < 5; i++ {
		go func(workerid int) {
			for t := range taskch {
				t.runnable()
			}
		}(i)
	}

	for i := 0; i < 100000; i++ {
		t := task{
			runnable: func() {
				lock.Lock()
				sum++
				lock.Unlock()
			},
		}
		taskch <- t
	}

	fmt.Println(sum)
}
