package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRocketChatAlerter_EmptyURL(t *testing.T) {
	_, err := NewRocketChatAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewRocketChatAlerter_ValidURL(t *testing.T) {
	a, err := NewRocketChatAlerter("https://chat.example.com/hooks/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestRocketChatAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a, err := NewRocketChatAlerter(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a.client = rewriteToServer(ts, a.client)

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/my-app/db", expiry); err != nil {
		t.Fatalf("unexpected Send error: %v", err)
	}

	if received["text"] == "" {
		t.Error("expected non-empty text field in payload")
	}
}

func TestRocketChatAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a, err := NewRocketChatAlerter(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a.client = rewriteToServer(ts, a.client)

	err = a.Send("secret/my-app/db", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestRocketChatAlerter_Send_PayloadHasAttachments(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a, err := NewRocketChatAlerter(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a.client = rewriteToServer(ts, a.client)

	_ = a.Send("secret/prod/token", time.Now().Add(2*time.Hour))

	attachments, ok := received["attachments"]
	if !ok {
		t.Fatal("expected attachments field in payload")
	}
	list, ok := attachments.([]interface{})
	if !ok || len(list) == 0 {
		t.Error("expected at least one attachment")
	}
}
