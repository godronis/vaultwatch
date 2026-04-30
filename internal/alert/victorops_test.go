package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewVictorOpsAlerter_EmptyURL(t *testing.T) {
	_, err := NewVictorOpsAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewVictorOpsAlerter_ValidURL(t *testing.T) {
	a, err := NewVictorOpsAlerter("https://alert.victorops.com/integrations/generic/123/alert/token/routing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestVictorOpsAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, err := NewVictorOpsAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server, a.client)

	if err := a.Send("secret/db/password", "expires in 2 days"); err != nil {
		t.Fatalf("Send returned unexpected error: %v", err)
	}

	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %v", received["message_type"])
	}
	if received["entity_id"] != "secret/db/password" {
		t.Errorf("expected entity_id 'secret/db/password', got %v", received["entity_id"])
	}
	if received["state_message"] != "expires in 2 days" {
		t.Errorf("expected state_message 'expires in 2 days', got %v", received["state_message"])
	}
}

func TestVictorOpsAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, err := NewVictorOpsAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server, a.client)

	if err := a.Send("secret/api/key", "expires soon"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestVictorOpsAlerter_Send_PayloadHasTimestamp(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, err := NewVictorOpsAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server, a.client)
	_ = a.Send("secret/cert", "expiring")

	if _, ok := received["timestamp"]; !ok {
		t.Error("expected timestamp field in payload")
	}
}
