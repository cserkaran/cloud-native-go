package retry

import (
	"context"
	"log"
	"time"
)

// Effector for retry logic. The function signature of the failing method,
// which needs to be retried, must take the form of effector.
// Effector is the function that interacts with the potentially failing service.
type Effector func(context.Context) (string, error)

// Retry function, which wraps the Effector function(the potentially failing method)
// and adds the retry logic.
// Retry function accepts Effector and returns a closure with the same function signature as Effector.
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
