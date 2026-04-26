package vault

import (
	"testing"
	"time"
)

func TestIsExpiringSoon_WithinWindow(t *testing.T) {
	expiresAt := time.Now().Add(10 * time.Minute)
	warnBefore := 30 * time.Minute

	if !isExpiringSoon(expiresAt, warnBefore) {
		t.Error("expected secret expiring in 10m to be flagged with 30m warn window")
	}
}

func TestIsExpiringSoon_OutsideWindow(t *testing.T) {
	expiresAt := time.Now().Add(2 * time.Hour)
	warnBefore := 30 * time.Minute

	if isExpiringSoon(expiresAt, warnBefore) {
		t.Error("expected secret expiring in 2h to NOT be flagged with 30m warn window")
	}
}

func TestIsExpiringSoon_AlreadyExpired(t *testing.T) {
	expiresAt := time.Now().Add(-5 * time.Minute)
	warnBefore := 30 * time.Minute

	if !isExpiringSoon(expiresAt, warnBefore) {
		t.Error("expected already-expired secret to be flagged")
	}
}

func TestMonitor_NoAlerts(t *testing.T) {
	// Monitor with an empty path list should return no alerts and no error.
	client := &Client{} // vc is nil; no calls will be made
	alerts, err := Monitor(client, []string{}, 30*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestMonitor_SkipsBadPaths(t *testing.T) {
	// A nil vc will cause GetSecretMeta to fail; Monitor should skip and continue.
	client := &Client{}
	alerts, err := Monitor(client, []string{"secret/data/missing"}, 30*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Errors are logged and skipped, so we expect 0 alerts (not a hard failure).
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts for bad path, got %d", len(alerts))
	}
}
