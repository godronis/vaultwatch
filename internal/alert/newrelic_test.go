package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewNewRelicAlerter_EmptyAccountID(t *testing.T) {
	_, err := NewNewRelicAlerter("", "some-api-key")
	if err == nil {
		t.Fatal("expected error for empty account ID, got nil")
	}
}

func TestNewNewRelicAlerter_EmptyAPIKey(t *testing.T) {
	_, err := NewNewRelicAlerter("123456", "")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
}

func TestNewNewRelicAlerter_ValidConfig(t *testing.T) {
	a, err := NewNewRelicAlerter("123456", "test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestNewRelicAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	var insertKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		insertKey = r.Header.Get("X-Insert-Key")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewNewRelicAlerter("123456", "test-key")
	a.endpoint = server.URL

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiresAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if insertKey != "test-key" {
		t.Errorf("expected X-Insert-Key 'test-key', got %q", insertKey)
	}
	if received["eventType"] != "VaultSecretExpiring" {
		t.Errorf("expected eventType 'VaultSecretExpiring', got %v", received["eventType"])
	}
	if received["secretPath"] != "secret/myapp/db" {
		t.Errorf("expected secretPath 'secret/myapp/db', got %v", received["secretPath"])
	}
}

func TestNewRelicAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	a, _ := NewNewRelicAlerter("123456", "bad-key")
	a.endpoint = server.URL

	err := a.Send("secret/myapp/db", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestNewRelicAlerter_Send_EventTypeField(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewNewRelicAlerter("123456", "test-key")
	a.endpoint = server.URL

	_ = a.Send("secret/infra/tls", time.Now().Add(48*time.Hour))

	if _, ok := received["expiresAt"]; !ok {
		t.Error("expected 'expiresAt' field in payload")
	}
	if _, ok := received["message"]; !ok {
		t.Error("expected 'message' field in payload")
	}
}
