package alert

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// RetryConfig holds configuration for HTTP retry behavior.
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

// doWithRetry executes the provided function with exponential backoff retries.
// It retries on HTTP 5xx responses or transport-level errors.
func doWithRetry(cfg RetryConfig, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		resp, err := fn()
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt, err)
		} else {
			lastErr = fmt.Errorf("attempt %d: server returned status %d", attempt, resp.StatusCode)
			resp.Body.Close()
		}

		if attempt < cfg.MaxAttempts {
			time.Sleep(delay)
			next := time.Duration(math.Min(
				float64(delay)*cfg.Multiplier,
				float64(cfg.MaxDelay),
			))
			delay = next
		}
	}

	return nil, fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}
