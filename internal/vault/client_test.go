package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_ValidAddress(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetSecretMeta_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetSecretMeta("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestGetSecretMeta_WithTTL(t *testing.T) {
	body := `{"data":{"ttl":3600}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	meta, err := client.GetSecretMeta("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.TTL.Seconds() != 3600 {
		t.Errorf("expected TTL 3600s, got %v", meta.TTL)
	}
	if meta.Path != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", meta.Path)
	}
}
