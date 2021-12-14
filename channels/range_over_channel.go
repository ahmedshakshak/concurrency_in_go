package main

import "fmt"

func main() {

	queue := make(chan string, 2)
	go func() {
		queue <- "one"
		queue <- "two"
		close(queue)
	}()

	for elem := range queue {
		fmt.Println(elem)
	}
}
