package circuit

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// The Circuit function, which interacts with the potentially failing service.
type Circuit func(context.Context) (string, error)

// New constructor to create the circuit function.
// Returns an anonymous function which fails intermittently.
// Failure condition is randomized. 
func New() (func(context.Context)(string,error)){
    return func(ctx context.Context) (string, error) {
		time.Sleep(3 * time.Second)
		random := rand.Intn(100)
		if random%3 == 0 {
			return "success", nil
		}
		return "", errors.New("error calling circuit logic")
	}
}
