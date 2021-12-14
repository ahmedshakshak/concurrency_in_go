package main

import (
	"fmt"
	"time"
)

func orDoneV1(channel chan interface{}) {
	for {
		val, ok := <-channel
		if ok == false {
			return // or maybe break from for
		}
		// Do something with val
		fmt.Println(val)
	}
}

func orDoneV2(done, channel chan interface{}) {
	for {
		select {
		case <-done:
			break
		case val, ok := <-channel:
			if ok == false {
				return // or maybe break from for
			}
			// Do something with val
			fmt.Println(val)
		}
	}
}

func main() {
	orDone := func(done, c <-chan int) <-chan int {
		valChannel := make(chan int)

		go func() {
			defer close(valChannel)

			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}

					select {
					case valChannel <- v:
					case <-done:
					}
				}
			}
		}()

		return valChannel
	}

	channel := make(chan int, 5)
	done := make(chan int)

	go func() {
		channel <- 1
		channel <- 2
		channel <- 3
		channel <- 4
		channel <- 5
		channel <- 6
	}()

	go func() {
		time.Sleep(1 * time.Second)
		done <- 0
	}()

	for val := range orDone(done, channel) {
		fmt.Printf("Hello from goroutine %v\n", val)
	}
}
