package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewPagerDutyAlerter_EmptyKey(t *testing.T) {
	_, err := NewPagerDutyAlerter("")
	if err == nil {
		t.Error("expected error for empty integration key")
	}
}

func TestNewPagerDutyAlerter_ValidKey(t *testing.T) {
	a, err := NewPagerDutyAlerter("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Error("expected non-nil alerter")
	}
}

func TestPagerDutyAlerter_Send_Success(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a, err := NewPagerDutyAlerter("my-routing-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client.Transport = rewriteTransport(ts.URL)

	expiry := time.Now().Add(3 * time.Hour)
	if err := a.Send("secret/prod/api", expiry); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if received.RoutingKey != "my-routing-key" {
		t.Errorf("expected routing key 'my-routing-key', got '%s'", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got '%s'", received.EventAction)
	}
	if !strings.Contains(received.Payload.Summary, "secret/prod/api") {
		t.Errorf("expected summary to contain path, got '%s'", received.Payload.Summary)
	}
}

func TestPagerDutyAlerter_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	a, err := NewPagerDutyAlerter("my-routing-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.client.Transport = rewriteTransport(ts.URL)

	expiry := time.Now().Add(1 * time.Hour)
	if err := a.Send("secret/prod/api", expiry); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

// rewriteTransport redirects all requests to the given test server URL.
type urlRewriter struct{ base string }

func rewriteTransport(base string) http.RoundTripper { return &urlRewriter{base: base} }

func (u *urlRewriter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(u.base, "http://")
	return http.DefaultTransport.RoundTrip(req)
}
