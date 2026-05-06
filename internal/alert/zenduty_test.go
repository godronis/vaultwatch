package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewZendutyAlerter_EmptyKey(t *testing.T) {
	_, err := NewZendutyAlerter("", "svc-123", "")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewZendutyAlerter_EmptyServiceID(t *testing.T) {
	_, err := NewZendutyAlerter("key-abc", "", "")
	if err == nil {
		t.Fatal("expected error for empty service ID")
	}
}

func TestNewZendutyAlerter_ValidConfig(t *testing.T) {
	a, err := NewZendutyAlerter("key-abc", "svc-123", "esc-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestZendutyAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token key-abc" {
			t.Errorf("missing or wrong Authorization header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	a, _ := NewZendutyAlerter("key-abc", "svc-123", "")
	a.client = rewriteToServer(ts)

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiresAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["service"] != "svc-123" {
		t.Errorf("expected service svc-123, got %v", received["service"])
	}
}

func TestZendutyAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	a, _ := NewZendutyAlerter("bad-key", "svc-123", "")
	a.client = rewriteToServer(ts)

	err := a.Send("secret/myapp/db", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestZendutyAlerter_Send_EscalationPolicyIncluded(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	a, _ := NewZendutyAlerter("key-abc", "svc-123", "esc-789")
	a.client = rewriteToServer(ts)

	_ = a.Send("secret/myapp/token", time.Now().Add(2*time.Hour))
	if received["escalation_policy"] != "esc-789" {
		t.Errorf("expected escalation_policy esc-789, got %v", received["escalation_policy"])
	}
}
