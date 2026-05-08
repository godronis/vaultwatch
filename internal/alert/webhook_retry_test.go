package alert

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestDoWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	cfg := DefaultRetryConfig()
	cfg.InitialDelay = time.Millisecond

	resp, err := doWithRetry(cfg, func() (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDoWithRetry_RetriesOnTransportError(t *testing.T) {
	calls := 0
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = time.Millisecond

	_, err := doWithRetry(cfg, func() (*http.Response, error) {
		calls++
		return nil, errors.New("connection refused")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDoWithRetry_RetriesOn500(t *testing.T) {
	calls := 0
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = time.Millisecond

	_, err := doWithRetry(cfg, func() (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody}, nil
	})

	if err == nil {
		t.Fatal("expected error after retries, got nil")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDoWithRetry_SucceedsAfterRetry(t *testing.T) {
	calls := 0
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = time.Millisecond

	resp, err := doWithRetry(cfg, func() (*http.Response, error) {
		calls++
		if calls < 3 {
			return nil, errors.New("transient error")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})

	if err != nil {
		t.Fatalf("expected success on third attempt, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDefaultRetryConfig_Values(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
