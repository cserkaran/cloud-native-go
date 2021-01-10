package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloud-native-go/retry"
)

var count int

func emulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

func main() {
	r := retry.Retry(emulateTransientError, 5, 2*time.Second)
	res, err := r(context.Background())
	fmt.Println(res, err)
}
