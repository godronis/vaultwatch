package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSlackAlerter_EmptyURL(t *testing.T) {
	_, err := NewSlackAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewSlackAlerter_ValidURL(t *testing.T) {
	a, err := NewSlackAlerter("https://hooks.slack.com/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.WebhookURL != "https://hooks.slack.com/test" {
		t.Errorf("expected webhook URL to be set")
	}
}

func TestSlackAlerter_Send_Success(t *testing.T) {
	var received SlackPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a, _ := NewSlackAlerter(ts.URL)
	expiry := time.Now().Add(24 * time.Hour)

	if err := a.Send("secret/myapp/db", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Text == "" {
		t.Error("expected non-empty text in slack payload")
	}
	if len(received.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(received.Attachments))
	}
	if received.Attachments[0].Title != "secret/myapp/db" {
		t.Errorf("unexpected attachment title: %s", received.Attachments[0].Title)
	}
}

func TestSlackAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a, _ := NewSlackAlerter(ts.URL)
	err := a.Send("secret/myapp/db", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}
