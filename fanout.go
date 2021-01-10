// Example for fanout pattern.

package main

import (
	"fmt"
	"sync"
)

// Split returns the number n of destination channels, 
// which will read from the single source and do some work,
// and are requested then to return the result via recieve.
func split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests = append(dests, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}

	return dests
}

func main() {
	source := make(chan int)
	dests := split(source, 5)

	go func() {
		for i := 1; i <= 10; i++ {
			source <- i
		}

		close(source)
	}()

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, ch := range dests {
		go func(i int, d <-chan int) {
			defer wg.Done()

			for val := range d {
				fmt.Printf("#%d got %d\n", i, val)
			}
		}(i, ch)
	}

	wg.Wait()
}
