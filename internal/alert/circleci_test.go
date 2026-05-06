package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewCircleCIAlerter_EmptyToken(t *testing.T) {
	_, err := NewCircleCIAlerter("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewCircleCIAlerter_ValidToken(t *testing.T) {
	a, err := NewCircleCIAlerter("mytoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestCircleCIAlerter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	a, _ := NewCircleCIAlerter("tok")
	a = rewriteToServer(server, a).(*CircleCIAlerter)

	if err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCircleCIAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a, _ := NewCircleCIAlerter("tok")
	a = rewriteToServer(server, a).(*CircleCIAlerter)

	if err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestCircleCIAlerter_Send_PayloadFields(t *testing.T) {
	var captured map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	a, _ := NewCircleCIAlerter("tok")
	a = rewriteToServer(server, a).(*CircleCIAlerter)

	expiry := time.Now().Add(48 * time.Hour)
	a.Send("secret/prod/key", expiry)

	if captured["name"] != "vault-secret-expiry" {
		t.Errorf("expected name 'vault-secret-expiry', got %v", captured["name"])
	}
	attrs, ok := captured["attributes"].(map[string]interface{})
	if !ok {
		t.Fatal("expected attributes map")
	}
	if attrs["secret_path"] != "secret/prod/key" {
		t.Errorf("unexpected secret_path: %v", attrs["secret_path"])
	}
}
