package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewOpsGenieAlerter_EmptyKey(t *testing.T) {
	_, err := NewOpsGenieAlerter("")
	if err == nil {
		t.Fatal("expected error for empty api key, got nil")
	}
}

func TestNewOpsGenieAlerter_ValidKey(t *testing.T) {
	a, err := NewOpsGenieAlerter("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestOpsGenieAlerter_Send_Success(t *testing.T) {
	var received opsGeniePayload
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a, err := NewOpsGenieAlerter("my-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.apiURL = ts.URL

	err = a.Send("secret/db/password", 48*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(authHeader, "GenieKey my-api-key") {
		t.Errorf("expected GenieKey auth header, got: %s", authHeader)
	}
	if !strings.Contains(received.Message, "secret/db/password") {
		t.Errorf("expected message to contain path, got: %s", received.Message)
	}
	if received.Details["path"] != "secret/db/password" {
		t.Errorf("expected details.path to be set, got: %s", received.Details["path"])
	}
}

func TestOpsGenieAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	a, err := NewOpsGenieAlerter("bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.apiURL = ts.URL

	err = a.Send("secret/api/token", 24*time.Hour)
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to mention status code 401, got: %v", err)
	}
}
