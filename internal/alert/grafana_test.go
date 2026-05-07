package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewGrafanaAlerter_EmptyURL(t *testing.T) {
	_, err := NewGrafanaAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewGrafanaAlerter_ValidURL(t *testing.T) {
	a, err := NewGrafanaAlerter("https://grafana.example.com/webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestGrafanaAlerter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewGrafanaAlerter(server.URL)
	a.client = rewriteToServer(server, a.client)

	err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGrafanaAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, _ := NewGrafanaAlerter(server.URL)
	a.client = rewriteToServer(server, a.client)

	err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestGrafanaAlerter_Send_PayloadFields(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewGrafanaAlerter(server.URL)
	a.client = rewriteToServer(server, a.client)

	_ = a.Send("secret/myapp/token", time.Now().Add(12*time.Hour))

	if _, ok := received["title"]; !ok {
		t.Error("expected 'title' field in payload")
	}
	if _, ok := received["state"]; !ok {
		t.Error("expected 'state' field in payload")
	}
	if received["state"] != "alerting" {
		t.Errorf("expected state 'alerting', got %v", received["state"])
	}
}
