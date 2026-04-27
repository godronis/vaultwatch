package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewTeamsAlerter_EmptyURL(t *testing.T) {
	_, err := NewTeamsAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewTeamsAlerter_ValidURL(t *testing.T) {
	a, err := NewTeamsAlerter("https://outlook.office.com/webhook/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestTeamsAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewTeamsAlerter(server.URL)
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := a.Send("secret/my-service/api-key", expiresAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["@type"] != "MessageCard" {
		t.Errorf("expected @type=MessageCard, got %v", received["@type"])
	}
	if !strings.Contains(received["summary"].(string), "secret/my-service/api-key") {
		t.Errorf("expected summary to contain path, got %v", received["summary"])
	}
}

func TestTeamsAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	a, _ := NewTeamsAlerter(server.URL)
	err := a.Send("secret/test", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "503") {
		t.Errorf("expected error to mention status 503, got: %v", err)
	}
}

func TestTeamsAlerter_Send_PayloadFields(t *testing.T) {
	var sections []interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &payload)
		sections, _ = payload["sections"].([]interface{})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewTeamsAlerter(server.URL)
	_ = a.Send("secret/db/password", time.Now().Add(2*time.Hour))

	if len(sections) == 0 {
		t.Fatal("expected at least one section in Teams payload")
	}
	section := sections[0].(map[string]interface{})
	facts, ok := section["facts"].([]interface{})
	if !ok || len(facts) < 3 {
		t.Errorf("expected at least 3 facts in section, got: %v", facts)
	}
}
