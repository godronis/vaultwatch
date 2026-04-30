package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewDiscordAlerter_EmptyURL(t *testing.T) {
	_, err := NewDiscordAlerter("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewDiscordAlerter_ValidURL(t *testing.T) {
	a, err := NewDiscordAlerter("https://discord.com/api/webhooks/123/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestDiscordAlerter_Send_Success(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	a, err := NewDiscordAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server)

	if err := a.Send("secret/myapp/db", "expires in 2 days"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload discordPayload
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if len(payload.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(payload.Embeds))
	}
	if !strings.Contains(payload.Embeds[0].Title, "secret/myapp/db") {
		t.Errorf("embed title missing path, got: %s", payload.Embeds[0].Title)
	}
	if !strings.Contains(payload.Embeds[0].Description, "expires in 2 days") {
		t.Errorf("embed description missing message, got: %s", payload.Embeds[0].Description)
	}
}

func TestDiscordAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, err := NewDiscordAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server)

	if err := a.Send("secret/myapp/db", "expires soon"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestDiscordAlerter_Send_EmbedColor(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	a, err := NewDiscordAlerter(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client = rewriteToServer(server)

	_ = a.Send("secret/x", "msg")

	var payload discordPayload
	_ = json.Unmarshal(received, &payload)
	if payload.Embeds[0].Color != 15158332 {
		t.Errorf("expected red color 15158332, got %d", payload.Embeds[0].Color)
	}
}
