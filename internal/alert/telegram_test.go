package alert

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewTelegramAlerter_EmptyToken(t *testing.T) {
	_, err := NewTelegramAlerter("", "12345")
	if err == nil {
		t.Fatal("expected error for empty bot token")
	}
}

func TestNewTelegramAlerter_EmptyChatID(t *testing.T) {
	_, err := NewTelegramAlerter("bot123", "")
	if err == nil {
		t.Fatal("expected error for empty chat ID")
	}
}

func TestNewTelegramAlerter_ValidConfig(t *testing.T) {
	a, err := NewTelegramAlerter("bot123:ABC", "-1001234567890")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

// newTestAlerter creates a TelegramAlerter wired to the given test server.
func newTestAlerter(t *testing.T, server *httptest.Server, botToken, chatID string) *TelegramAlerter {
	t.Helper()
	a := &TelegramAlerter{
		botToken: botToken,
		chatID:   chatID,
		client:   server.Client(),
	}
	a.client.Transport = rewriteToServer(server.URL)
	return a
}

func TestTelegramAlerter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a := newTestAlerter(t, server, "testtoken", "12345")

	secret := vault.SecretMeta{
		Path:       "secret/db/password",
		Expiration: time.Now().Add(24 * time.Hour),
	}
	if err := a.Send(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTelegramAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	a := newTestAlerter(t, server, "badtoken", "12345")

	secret := vault.SecretMeta{
		Path:       "secret/api/key",
		Expiration: time.Now().Add(1 * time.Hour),
	}
	err := a.Send(secret)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
