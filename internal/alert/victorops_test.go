package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewVictorOpsAlerter_EmptyURL(t *testing.T) {
	_, err := NewVictorOpsAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewVictorOpsAlerter_ValidURL(t *testing.T) {
	a, err := NewVictorOpsAlerter("https://alert.victorops.com/integrations/generic/12345/alert/token/routing")
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
	a.client = rewriteToServer(server)

	p := Payload{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := a.Send(p); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %v", received["message_type"])
	}
	if received["entity_id"] != "secret/myapp/db" {
		t.Errorf("expected entity_id 'secret/myapp/db', got %v", received["entity_id"])
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
	a.client = rewriteToServer(server)

	p := Payload{Path: "secret/api/key", ExpiresAt: time.Now().Add(time.Hour)}
	if err := a.Send(p); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}
