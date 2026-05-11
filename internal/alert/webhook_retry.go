package alert

import (
	"fmt"
	"net/http"
	"time"
)

// RetryConfig defines the parameters for retry behavior on webhook delivery.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the first retry.
	InitialDelay time.Duration
	// MaxDelay caps the exponential backoff delay.
	MaxDelay time.Duration
	// Multiplier is the factor by which the delay grows on each retry.
	Multiplier float64
}

// DefaultRetryConfig returns a sensible default retry configuration:
// 3 attempts, starting at 500ms, capped at 10s, doubling each time.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// doWithRetry executes fn up to cfg.MaxAttempts times, retrying on transport
// errors or HTTP 5xx responses. It returns the last response and error.
func doWithRetry(cfg RetryConfig, fn func() (*http.Response, error)) (*http.Response, error) {
	var (
		resp  *http.Response
		err   error
		delay = cfg.InitialDelay
	)

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		resp, err = fn()

		// Success: no transport error and not a server-side error.
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Don't sleep after the final attempt.
		if attempt == cfg.MaxAttempts {
			break
		}

		// Drain and close the body before retrying to allow connection reuse.
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		time.Sleep(delay)

		// Exponential backoff, capped at MaxDelay.
		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	if err != nil {
		return nil, fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, err)
	}
	return resp, nil
}
