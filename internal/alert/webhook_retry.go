package alert

import (
	"fmt"
	"net/http"
	"time"
)

// RetryConfig holds configuration for retry behaviour on HTTP alerters.
type RetryConfig struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// doWithRetry executes fn, retrying on transport errors or 5xx responses.
// It respects the provided RetryConfig for backoff behaviour.
func doWithRetry(cfg RetryConfig, fn func() (*http.Response, error)) (*http.Response, error) {
	delay := cfg.InitialDelay
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		resp, err := fn()
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt, err)
		} else {
			lastErr = fmt.Errorf("attempt %d: server returned %d", attempt, resp.StatusCode)
			resp.Body.Close()
		}
		if attempt < cfg.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}
	return nil, fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}
