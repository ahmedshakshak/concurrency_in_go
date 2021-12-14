package main

import (
	"fmt"
	"time"
)

func orDone(done, c <-chan int) <-chan int {
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

func main() {
	tee := func(done, in <-chan int) (_, _ <-chan int) {
		out1 := make(chan int)
		out2 := make(chan int)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
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

	out1, out2 := tee(done, channel)

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
