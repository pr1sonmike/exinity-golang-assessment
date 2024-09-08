package utils

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

func Retry(ctx context.Context, maxRetries int, initialBackoff time.Duration, operation func() error) error {
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		sleepDuration := backoff + time.Duration(rand.Int63n(int64(backoff/2)))
		select {
		case <-time.After(sleepDuration):
			backoff *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.New("operation failed after max retries")
}
