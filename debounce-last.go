package main

import (
	"context"
	"fmt"
	"time"
	"github.com/cloud-native-go/circuit"
)

//DebounceLast function, the Function last implementation of debounce.
// It wraps around the circuit logic, calling it the only the last time time when a series
// of calls are made in a cluster during a time duration.
// Employs a time ticker, to determine if enough time has passed since the function
// was last called.
func DebounceLast(circuit circuit.Circuit, d time.Duration) circuit.Circuit {
	var threshold time.Time = time.Now()
	var ticker *time.Ticker
	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		threshold = time.Now().Add(d)

		if ticker == nil {
			ticker = time.NewTicker(time.Millisecond * 100)
			tickerc := ticker.C

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-tickerc:
						if threshold.Before(time.Now()) {
							result, err = circuit(ctx)
							ticker.Stop()
							ticker = nil
							break
						}
					case <-ctx.Done():
						result, err = "", ctx.Err()
						break
					}
				}
			}()
		}

		return result, err
	}
}

func main() {

	ckt := circuit.New()
	ctx := context.Background()
	debounceLast := DebounceLast(ckt, 5 * time.Second)
	for {

		res, err := debounceLast(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}
}
