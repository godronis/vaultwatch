package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewVictorOpsV2Alerter_EmptyEndpoint(t *testing.T) {
	_, err := NewVictorOpsV2Alerter("", "routing-key")
	if err == nil {
		t.Fatal("expected error for empty endpoint")
	}
}

func TestNewVictorOpsV2Alerter_EmptyRoutingKey(t *testing.T) {
	_, err := NewVictorOpsV2Alerter("https://alert.victorops.com/integrations/generic", "")
	if err == nil {
		t.Fatal("expected error for empty routing key")
	}
}

func TestNewVictorOpsV2Alerter_ValidConfig(t *testing.T) {
	a, err := NewVictorOpsV2Alerter("https://alert.victorops.com/integrations/generic", "my-routing-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestVictorOpsV2Alerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewVictorOpsV2Alerter(server.URL, "my-key")
	a.client = rewriteToServer(server, a.client)

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVictorOpsV2Alerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, _ := NewVictorOpsV2Alerter(server.URL, "my-key")
	a.client = rewriteToServer(server, a.client)

	if err := a.Send("secret/myapp/db", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestVictorOpsV2Alerter_Send_PayloadFields(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewVictorOpsV2Alerter(server.URL, "my-key")
	a.client = rewriteToServer(server, a.client)

	a.Send("secret/myapp/db", time.Now().Add(time.Hour))

	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %v", received["message_type"])
	}
	if received["monitoring_tool"] != "vaultwatch" {
		t.Errorf("expected monitoring_tool vaultwatch, got %v", received["monitoring_tool"])
	}
}
