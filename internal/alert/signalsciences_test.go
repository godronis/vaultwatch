package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSignalSciencesAlerter_EmptyCorpName(t *testing.T) {
	_, err := NewSignalSciencesAlerter("", "mysite", "user@example.com", "token")
	if err == nil {
		t.Fatal("expected error for empty corp name")
	}
}

func TestNewSignalSciencesAlerter_EmptySiteName(t *testing.T) {
	_, err := NewSignalSciencesAlerter("mycorp", "", "user@example.com", "token")
	if err == nil {
		t.Fatal("expected error for empty site name")
	}
}

func TestNewSignalSciencesAlerter_EmptyEmail(t *testing.T) {
	_, err := NewSignalSciencesAlerter("mycorp", "mysite", "", "token")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestNewSignalSciencesAlerter_EmptyToken(t *testing.T) {
	_, err := NewSignalSciencesAlerter("mycorp", "mysite", "user@example.com", "")
	if err == nil {
		t.Fatal("expected error for empty API token")
	}
}

func TestNewSignalSciencesAlerter_ValidConfig(t *testing.T) {
	a, err := NewSignalSciencesAlerter("mycorp", "mysite", "user@example.com", "secret-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestSignalSciencesAlerter_Send_Success(t *testing.T) {
	var gotEmail, gotToken string
	var body map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotEmail = r.Header.Get("x-api-user")
		gotToken = r.Header.Get("x-api-token")
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a, _ := NewSignalSciencesAlerter("mycorp", "mysite", "user@example.com", "my-token")
	a.httpClient = rewriteToServer(server)

	expiry := time.Now().Add(24 * time.Hour)
	if err := a.Send("secret/myapp/db", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotEmail != "user@example.com" {
		t.Errorf("expected email header, got %q", gotEmail)
	}
	if gotToken != "my-token" {
		t.Errorf("expected token header, got %q", gotToken)
	}
	if body["message"] == nil {
		t.Error("expected message field in payload")
	}
}

func TestSignalSciencesAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	a, _ := NewSignalSciencesAlerter("mycorp", "mysite", "user@example.com", "bad-token")
	a.httpClient = rewriteToServer(server)

	if err := a.Send("secret/myapp/db", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
