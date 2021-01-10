// Example of Debounce first Cloud native pattern.
package main

import (
	"context"
	"fmt"
	"github.com/cloud-native-go/circuit"
	"time"
)

// DebounceFirst function, the function first implementation of debounce.
// It wraps around the circuit logic, calling it the only the first time when a series
// of calls are made in a cluster during a time duration.
func DebounceFirst(circuit circuit.Circuit, d time.Duration) circuit.Circuit {
	var threshold time.Time
	var cResult string
	var cError error

	return func(ctx context.Context) (string, error) {
		if threshold.Before(time.Now()) {
			cResult, cError = circuit(ctx)
		}

		threshold = time.Now().Add(d)
		return cResult, cError
	}
}

func main() {
	ckt := circuit.New()
	ctx := context.Background()
	debounceFirst := DebounceFirst(ckt, 5 * time.Second)
	for {

		res, err := debounceFirst(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}
}
