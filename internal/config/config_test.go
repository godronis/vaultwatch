package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "vaultwatch-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: http://127.0.0.1:8200
  token: s.test
secrets:
  - secret/myapp/db
warn_before: 48h
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected address: %s", cfg.Vault.Address)
	}
	if cfg.WarnBefore != 48*time.Hour {
		t.Errorf("unexpected warn_before: %v", cfg.WarnBefore)
	}
}

func TestLoad_DefaultWarnBefore(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: http://127.0.0.1:8200
  token: s.test
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WarnBefore != 72*time.Hour {
		t.Errorf("expected default 72h, got %v", cfg.WarnBefore)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  token: s.test
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: http://127.0.0.1:8200
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestLoad_TeamsConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: http://127.0.0.1:8200
  token: s.test
alerts:
  teams:
    webhook_url: https://outlook.office.com/webhook/abc123
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Alerts.Teams == nil {
		t.Fatal("expected Teams config to be set")
	}
	if cfg.Alerts.Teams.WebhookURL != "https://outlook.office.com/webhook/abc123" {
		t.Errorf("unexpected Teams webhook URL: %s", cfg.Alerts.Teams.WebhookURL)
	}
}
