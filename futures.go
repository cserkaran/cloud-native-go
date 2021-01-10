package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// the interface that is recieved by consumer to retrieve the eventual result.
type future interface {
	Result() (string, error)
}

// Satisfies the future interface, includes an attached method that contains 
// the result accessing logic.
type innerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	res   string
	err   error
	resCh <-chan string
	errCh <-chan error
}

// The results accessing logic, satisfies future interface.
// Does a recieve on each channel of the final result and 
// executes it once.
func (f *innerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}

// A wrapper function around some function to be asynchronously executed,
// provides future.
func slowFunction(ctx context.Context) future {
	resch := make(chan string)
	errch := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 2):
			resch <- "I slept for 2 seconds"
			errch <- nil
		case <-ctx.Done():
			resch <- ""
			errch <- ctx.Err()
		}
	}()

	return &innerFuture{resCh: resch, errCh: errch}
}

func main() {
	ctx := context.Background()
	future := slowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(res)
	}
}
