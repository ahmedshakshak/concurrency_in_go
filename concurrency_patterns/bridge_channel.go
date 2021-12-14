package main

import (
	"fmt"
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
	bridge := func(done <-chan int, chanStream <-chan <-chan int) <-chan int {
		valStream := make(chan int)
		go func() {
			defer close(valStream)
			for {
				var stream <-chan int
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	genVals := func() <-chan <-chan int {
		chanStream := make(chan (<-chan int))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan int, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
