package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSplunkAlerter_EmptyURL(t *testing.T) {
	_, err := NewSplunkAlerter("", "token", "vaultwatch")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewSplunkAlerter_EmptyToken(t *testing.T) {
	_, err := NewSplunkAlerter("http://splunk.example.com", "", "vaultwatch")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewSplunkAlerter_DefaultSource(t *testing.T) {
	a, err := NewSplunkAlerter("http://splunk.example.com", "mytoken", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.source != "vaultwatch" {
		t.Errorf("expected default source 'vaultwatch', got %q", a.source)
	}
}

func TestNewSplunkAlerter_ValidConfig(t *testing.T) {
	a, err := NewSplunkAlerter("http://splunk.example.com", "mytoken", "mysource")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.url != "http://splunk.example.com" {
		t.Errorf("unexpected url: %s", a.url)
	}
}

func TestSplunkAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, err := NewSplunkAlerter(server.URL, "test-token", "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiry); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if !strings.HasPrefix(authHeader, "Splunk ") {
		t.Errorf("expected Authorization header to start with 'Splunk ', got %q", authHeader)
	}
	if received["source"] != "vaultwatch" {
		t.Errorf("unexpected source: %v", received["source"])
	}
}

func TestSplunkAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	a, err := NewSplunkAlerter(server.URL, "bad-token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := a.Send("secret/myapp/db", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSplunkAlerter_Send_EventContainsPath(t *testing.T) {
	var received splunkEvent

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewSplunkAlerter(server.URL, "tok", "vaultwatch")
	expiry := time.Now().Add(2 * time.Hour)
	a.Send("secret/prod/api", expiry)

	path, _ := received.Event["path"].(string)
	if path != "secret/prod/api" {
		t.Errorf("expected path 'secret/prod/api' in event, got %q", path)
	}
}
