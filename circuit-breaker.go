// Example to demonstrate circuit breaker cloud native pattern.
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// The Circuit function, which interacts with the potenitally failing service.
type Circuit func(context.Context) (string, error)

// Breaker function, A closure with same function signature as Circuit. It adds extra error handling
// logic to the Circuit function, also adds exponential back off in case service 
// is continuosly failing.
func Breaker(circuit Circuit, failureThreshold uint64) Circuit {

	var lastStateSuccessul = true
	var consecutiveFailures uint64 = 0
	var lastAttempt time.Time = time.Now()

	return func(ctx context.Context) (string, error) {

		if consecutiveFailures >= failureThreshold {
			backOffLevel := consecutiveFailures - failureThreshold
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << backOffLevel)
			if !time.Now().After(shouldRetryAt) {
				return "", errors.New("circuit open -- service unreachable")
			}
		}

		lastAttempt = time.Now()
		response, err := circuit(ctx)

		if err != nil {
			if !lastStateSuccessul {
				consecutiveFailures++
			}
			lastStateSuccessul = false

			return response, err
		}

		lastStateSuccessul = true
		consecutiveFailures = 0
		return response, nil
	}

}

func main() {

	circuit := func(ctx context.Context) (string, error) {
		time.Sleep(3 * time.Second)
		random := rand.Intn(100)
		if random%3 == 0 {
			return "success", nil
		}
		return "", errors.New("error calling circuit logic")

	}

	ctx := context.Background()
	breaker := Breaker(circuit, 4)
	for {

		res, err := breaker(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}

}
