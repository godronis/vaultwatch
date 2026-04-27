package alert

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
)

func TestNewEmailAlerter_MissingHost(t *testing.T) {
	cfg := config.EmailConfig{
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	}
	_, err := NewEmailAlerter(cfg)
	if err == nil {
		t.Fatal("expected error for missing host, got nil")
	}
}

func TestNewEmailAlerter_MissingFrom(t *testing.T) {
	cfg := config.EmailConfig{
		Host: "smtp.example.com",
		To:   []string{"ops@example.com"},
	}
	_, err := NewEmailAlerter(cfg)
	if err == nil {
		t.Fatal("expected error for missing from address, got nil")
	}
}

func TestNewEmailAlerter_MissingRecipients(t *testing.T) {
	cfg := config.EmailConfig{
		Host: "smtp.example.com",
		From: "vault@example.com",
	}
	_, err := NewEmailAlerter(cfg)
	if err == nil {
		t.Fatal("expected error for missing recipients, got nil")
	}
}

func TestNewEmailAlerter_ValidConfig(t *testing.T) {
	cfg := config.EmailConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	}
	alerter, err := NewEmailAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alerter.host != "smtp.example.com" {
		t.Errorf("expected host smtp.example.com, got %s", alerter.host)
	}
	if alerter.port != 587 {
		t.Errorf("expected port 587, got %d", alerter.port)
	}
}

func TestNewEmailAlerter_DefaultPort(t *testing.T) {
	cfg := config.EmailConfig{
		Host: "smtp.example.com",
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	}
	alerter, err := NewEmailAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alerter.port != 587 {
		t.Errorf("expected default port 587, got %d", alerter.port)
	}
}

func TestBuildEmailBody_ContainsPath(t *testing.T) {
	payload := map[string]string{
		"expires_at": "2024-12-31T00:00:00Z",
		"ttl":        "72h",
	}
	body := buildEmailBody("secret/myapp/db", payload)
	if !contains(body, "secret/myapp/db") {
		t.Errorf("expected body to contain secret path, got: %s", body)
	}
	if !contains(body, "expires_at") {
		t.Errorf("expected body to contain expires_at field, got: %s", body)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
