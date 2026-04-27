package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildPayload_ContainsExpectedFields(t *testing.T) {
	expiry := time.Now().Add(24 * time.Hour)
	payload := BuildPayload("secret/myapp/db", expiry)

	if payload["path"] != "secret/myapp/db" {
		t.Errorf("expected path 'secret/myapp/db', got %v", payload["path"])
	}
	if _, ok := payload["expires_at"]; !ok {
		t.Error("expected 'expires_at' key in payload")
	}
	if _, ok := payload["message"]; !ok {
		t.Error("expected 'message' key in payload")
	}
}

func TestNewWebhookAlerter_EmptyURL(t *testing.T) {
	_, err := NewWebhookAlerter("", nil)
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestWebhookAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	alerter, err := NewWebhookAlerter(ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expiry := time.Now().Add(2 * time.Hour)
	if err := alerter.Send("secret/myapp/db", expiry); err != nil {
		t.Errorf("unexpected send error: %v", err)
	}
	if received["path"] != "secret/myapp/db" {
		t.Errorf("expected path in payload, got %v", received["path"])
	}
}

func TestWebhookAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	alerter, err := NewWebhookAlerter(ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expiry := time.Now().Add(2 * time.Hour)
	if err := alerter.Send("secret/myapp/db", expiry); err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestWebhookAlerter_Send_WithHeaders(t *testing.T) {
	var capturedHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("X-Custom-Token")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	headers := map[string]string{"X-Custom-Token": "abc123"}
	alerter, err := NewWebhookAlerter(ts.URL, headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expiry := time.Now().Add(2 * time.Hour)
	if err := alerter.Send("secret/path", expiry); err != nil {
		t.Errorf("unexpected send error: %v", err)
	}
	if capturedHeader != "abc123" {
		t.Errorf("expected header 'abc123', got '%s'", capturedHeader)
	}
}
