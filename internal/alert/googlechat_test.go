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

func TestNewGoogleChatAlerter_EmptyURL(t *testing.T) {
	_, err := NewGoogleChatAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewGoogleChatAlerter_ValidURL(t *testing.T) {
	a, err := NewGoogleChatAlerter("https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestGoogleChatAlerter_Send_Success(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a, _ := NewGoogleChatAlerter(ts.URL)
	a.client = rewriteToServer(ts)

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/my-app/api-key", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, ok := received["text"].(string)
	if !ok || !strings.Contains(text, "secret/my-app/api-key") {
		t.Errorf("expected payload text to contain path, got: %v", text)
	}
	if !strings.Contains(text, "VaultWatch Alert") {
		t.Errorf("expected payload text to contain 'VaultWatch Alert', got: %v", text)
	}
}

func TestGoogleChatAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	a, _ := NewGoogleChatAlerter(ts.URL)
	a.client = rewriteToServer(ts)

	err := a.Send("secret/my-app/api-key", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected error to mention status code 403, got: %v", err)
	}
}
