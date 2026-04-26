package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `
vault:
  address: "https://vault.example.com"
  token: "s.abc123"
  paths:
    - "secret/myapp"
  warn_before: 168h
alerts:
  webhooks:
    - name: slack
      url: "https://hooks.slack.com/services/XXX"
      timeout: 5s
schedule: "@hourly"
`
	path := writeTempConfig(t, raw)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if cfg.Vault.WarnBefore != 168*time.Hour {
		t.Errorf("expected 168h warn_before, got %v", cfg.Vault.WarnBefore)
	}
	if len(cfg.Alerts.Webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(cfg.Alerts.Webhooks))
	}
	if cfg.Alerts.Webhooks[0].Name != "slack" {
		t.Errorf("expected webhook name 'slack', got %q", cfg.Alerts.Webhooks[0].Name)
	}
}

func TestLoad_DefaultWarnBefore(t *testing.T) {
	raw := `
vault:
  address: "https://vault.example.com"
  token: "s.abc123"
  paths:
    - "secret/myapp"
`
	path := writeTempConfig(t, raw)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.WarnBefore != 7*24*time.Hour {
		t.Errorf("expected default 7d warn_before, got %v", cfg.Vault.WarnBefore)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	raw := `
vault:
  token: "s.abc123"
  paths:
    - "secret/myapp"
`
	path := writeTempConfig(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault.address, got nil")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	raw := `
vault:
  address: "https://vault.example.com"
  paths:
    - "secret/myapp"
`
	path := writeTempConfig(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault.token, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
