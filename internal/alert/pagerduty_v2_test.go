package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPagerDutyV2Alerter_EmptyKey(t *testing.T) {
	_, err := NewPagerDutyV2Alerter("")
	if err == nil {
		t.Fatal("expected error for empty integration key")
	}
}

func TestNewPagerDutyV2Alerter_ValidKey(t *testing.T) {
	a, err := NewPagerDutyV2Alerter("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestPagerDutyV2Alerter_Send_Success(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	a, _ := NewPagerDutyV2Alerter("key-abc")
	a.client = rewriteToServer(server)

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/my-app", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["routing_key"] != "key-abc" {
		t.Errorf("expected routing_key 'key-abc', got %v", received["routing_key"])
	}
	if received["event_action"] != "trigger" {
		t.Errorf("expected event_action 'trigger', got %v", received["event_action"])
	}
}

func TestPagerDutyV2Alerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	a, _ := NewPagerDutyV2Alerter("key-abc")
	a.client = rewriteToServer(server)

	if err := a.Send("secret/my-app", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyV2Alerter_Send_PayloadFields(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	a, _ := NewPagerDutyV2Alerter("key-xyz")
	a.client = rewriteToServer(server)

	expiry := time.Now().Add(12 * time.Hour)
	_ = a.Send("secret/db/password", expiry)

	inner, ok := received["payload"].(map[string]interface{})
	if !ok {
		t.Fatal("expected payload field to be a map")
	}
	if inner["severity"] != "warning" {
		t.Errorf("expected severity 'warning', got %v", inner["severity"])
	}
	if inner["source"] != "vaultwatch" {
		t.Errorf("expected source 'vaultwatch', got %v", inner["source"])
	}
}
