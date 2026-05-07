package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewStatusPageAlerter_EmptyAPIKey(t *testing.T) {
	_, err := NewStatusPageAlerter("", "page123", "comp456")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewStatusPageAlerter_EmptyPageID(t *testing.T) {
	_, err := NewStatusPageAlerter("key", "", "comp456")
	if err == nil {
		t.Fatal("expected error for empty page ID")
	}
}

func TestNewStatusPageAlerter_EmptyComponentID(t *testing.T) {
	_, err := NewStatusPageAlerter("key", "page123", "")
	if err == nil {
		t.Fatal("expected error for empty component ID")
	}
}

func TestNewStatusPageAlerter_ValidConfig(t *testing.T) {
	a, err := NewStatusPageAlerter("key", "page123", "comp456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestStatusPageAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer svr.Close()

	a, _ := NewStatusPageAlerter("mykey", "mypage", "mycomp")
	a.client = rewriteToServer(svr)

	err := a.Send("secret/my-cert", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	incident, ok := received["incident"].(map[string]interface{})
	if !ok {
		t.Fatal("expected incident key in payload")
	}
	if incident["status"] != "investigating" {
		t.Errorf("expected status 'investigating', got %v", incident["status"])
	}
}

func TestStatusPageAlerter_Send_NonOKStatus(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	a, _ := NewStatusPageAlerter("badkey", "page123", "comp456")
	a.client = rewriteToServer(svr)

	err := a.Send("secret/my-cert", time.Now().Add(24*time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
