package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewMattermostAlerter_EmptyURL(t *testing.T) {
	_, err := NewMattermostAlerter("", "", "")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewMattermostAlerter_ValidURL(t *testing.T) {
	a, err := NewMattermostAlerter("https://mattermost.example.com/hooks/abc", "#alerts", "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestMattermostAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, err := NewMattermostAlerter(server.URL, "#ops", "bot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiry); err != nil {
		t.Fatalf("Send returned unexpected error: %v", err)
	}

	text, _ := received["text"].(string)
	if !strings.Contains(text, "secret/myapp/db") {
		t.Errorf("expected payload text to contain secret path, got: %s", text)
	}
	if received["channel"] != "#ops" {
		t.Errorf("expected channel '#ops', got: %v", received["channel"])
	}
	if received["username"] != "bot" {
		t.Errorf("expected username 'bot', got: %v", received["username"])
	}
}

func TestMattermostAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, err := NewMattermostAlerter(server.URL, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = a.Send("secret/myapp/token", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code 500, got: %v", err)
	}
}

func TestMattermostAlerter_Send_OmitsEmptyFields(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, err := NewMattermostAlerter(server.URL, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a.Send("secret/path", time.Now().Add(time.Hour))

	if _, ok := received["channel"]; ok {
		t.Error("expected 'channel' to be omitted when empty")
	}
	if _, ok := received["username"]; ok {
		t.Error("expected 'username' to be omitted when empty")
	}
}
