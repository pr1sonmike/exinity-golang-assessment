package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_SuccessFirstAttempt(t *testing.T) {
	operation := func() error {
		return nil
	}

	err := Retry(context.Background(), 3, 100*time.Millisecond, operation)
	assert.NoError(t, err)
}

func TestRetry_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := Retry(context.Background(), 5, 100*time.Millisecond, operation)
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestRetry_FailureAfterMaxRetries(t *testing.T) {
	operation := func() error {
		return errors.New("permanent error")
	}

	err := Retry(context.Background(), 3, 100*time.Millisecond, operation)
	assert.Error(t, err)
	assert.Equal(t, "operation failed after max retries", err.Error())
}

func TestRetry_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	operation := func() error {
		return errors.New("temporary error")
	}

	err := Retry(ctx, 3, 100*time.Millisecond, operation)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}
