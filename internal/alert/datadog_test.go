package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewDatadogAlerter_EmptyKey(t *testing.T) {
	_, err := NewDatadogAlerter("")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
}

func TestNewDatadogAlerter_ValidKey(t *testing.T) {
	a, err := NewDatadogAlerter("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestDatadogAlerter_Send_Success(t *testing.T) {
	var received datadogPayload
	var apiKeyHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyHeader = r.Header.Get("DD-API-KEY")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a, _ := NewDatadogAlerter("my-api-key")
	a.baseURL = ts.URL

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiresAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if apiKeyHeader != "my-api-key" {
		t.Errorf("expected DD-API-KEY header 'my-api-key', got '%s'", apiKeyHeader)
	}
	if received.AlertType != "warning" {
		t.Errorf("expected alert_type 'warning', got '%s'", received.AlertType)
	}
	if received.Title != "Vault Secret Expiring Soon" {
		t.Errorf("unexpected title: %s", received.Title)
	}
}

func TestDatadogAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	a, _ := NewDatadogAlerter("bad-key")
	a.baseURL = ts.URL

	err := a.Send("secret/myapp/db", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestDatadogAlerter_Send_TagsContainPath(t *testing.T) {
	var received datadogPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a, _ := NewDatadogAlerter("key")
	a.baseURL = ts.URL
	a.Send("secret/myapp/token", time.Now().Add(time.Hour))

	found := false
	for _, tag := range received.Tags {
		if tag == "secret_path:secret/myapp/token" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected tag 'secret_path:secret/myapp/token' in tags %v", received.Tags)
	}
}
