package alert

import (
	"errors"
	"testing"
	"time"
)

// mockAlerter records calls and optionally returns an error.
type mockAlerter struct {
	called    bool
	path      string
	expiresAt time.Time
	errToReturn error
}

func (m *mockAlerter) Send(path string, expiresAt time.Time) error {
	m.called = true
	m.path = path
	m.expiresAt = expiresAt
	return m.errToReturn
}

func TestMultiAlerter_CallsAll(t *testing.T) {
	a1 := &mockAlerter{}
	a2 := &mockAlerter{}
	ma := NewMultiAlerter(a1, a2)

	expiry := time.Now().Add(time.Hour)
	if err := ma.Send("secret/path", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !a1.called || !a2.called {
		t.Error("expected both alerters to be called")
	}
	if a1.path != "secret/path" || a2.path != "secret/path" {
		t.Error("expected correct path forwarded to alerters")
	}
}

func TestMultiAlerter_ReturnsErrorsFromAll(t *testing.T) {
	a1 := &mockAlerter{errToReturn: errors.New("a1 failed")}
	a2 := &mockAlerter{errToReturn: errors.New("a2 failed")}
	ma := NewMultiAlerter(a1, a2)

	err := ma.Send("secret/path", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !a1.called || !a2.called {
		t.Error("expected both alerters to be called even when errors occur")
	}
}

func TestMultiAlerter_PartialFailure(t *testing.T) {
	a1 := &mockAlerter{}
	a2 := &mockAlerter{errToReturn: errors.New("a2 failed")}
	ma := NewMultiAlerter(a1, a2)

	err := ma.Send("secret/path", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error when one alerter fails")
	}
	if !a1.called {
		t.Error("expected a1 to be called despite a2 failure")
	}
}
