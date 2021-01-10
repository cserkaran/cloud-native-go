package throttle

import (
	"context"
	"time"
)

// Effector function, the method which interacts with the service.
// which needs to be rate limited by Throttling.
type Effector func(context.Context) (string, error)

// Throttle basic token bucket algoritm implementation, that uses the 
// "replay" strategy. Wraps the effector function in a closure that contains
// rate limiting logic.
// The bucket is initially allocated max tokens, each time the closure is trigerred
// it checks whether it has any remaining tokens. If yes, it decrements the token count
// by one and triggers effector. If not, last recorded result is replayed.
// Tokens are added at a rate of refill tokens every duration d.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var ticker *time.Ticker = nil
	var tokens uint = max

	var lastReturnString string
	var lastReturnError error

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		if ticker == nil {
			ticker = time.NewTicker(d)
			defer ticker.Stop()

			go func() {
				for {
					select {
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
                        }
                        tokens = t
					case <-ctx.Done():
						ticker.Stop()
						break
					}
				}
			}()
		}
		if tokens > 0 {
			tokens--
			lastReturnString, lastReturnError = e(ctx)
		}

		return lastReturnString, lastReturnError
	}
}
