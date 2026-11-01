// Package agent provides retry
package agent

import (
	"time"
)

// Retry retries function
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}

// RetryWithBackoff retries with backoff
func RetryWithBackoff(attempts int, initialDelay time.Duration, fn func() error) error {
	delay := initialDelay
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2
	}
	return err
}
