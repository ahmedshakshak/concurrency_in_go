package main

import (
	"fmt"
	"sync"
)

func main() {
	var data int
	var mutex sync.Mutex

	go func() {
		mutex.Lock()
		data++
		mutex.Unlock()
	}()

	mutex.Lock()
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	}
	mutex.Unlock()
}
