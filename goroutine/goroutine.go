package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("Hello from goroutine 1")
	}()

	fmt.Println("Hello from main goroutine")

	time.Sleep(1 * time.Second)
}
