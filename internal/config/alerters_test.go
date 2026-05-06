package config

import (
	"testing"
)

func TestAlertConfig_TelegramFields(t *testing.T) {
	cfg := AlertConfig{
		TelegramToken:  "bot123",
		TelegramChatID: "-100456",
	}
	if cfg.TelegramToken != "bot123" {
		t.Errorf("unexpected TelegramToken: %s", cfg.TelegramToken)
	}
	if cfg.TelegramChatID != "-100456" {
		t.Errorf("unexpected TelegramChatID: %s", cfg.TelegramChatID)
	}
}

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	cfg := AlertConfig{}
	if cfg.WebhookURL != "" {
		t.Error("expected empty WebhookURL")
	}
	if cfg.SlackURL != "" {
		t.Error("expected empty SlackURL")
	}
	if cfg.PagerDutyKey != "" {
		t.Error("expected empty PagerDutyKey")
	}
	if cfg.SplunkURL != "" {
		t.Error("expected empty SplunkURL")
	}
	if cfg.SplunkToken != "" {
		t.Error("expected empty SplunkToken")
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	cfg := AlertConfig{
		SlackURL:     "https://hooks.slack.com/xxx",
		PagerDutyKey: "pdkey",
		SplunkURL:    "https://splunk.example.com:8088/services/collector/event",
		SplunkToken:  "splunk-hec-token",
		SplunkSource: "vaultwatch-prod",
	}
	if cfg.SlackURL == "" {
		t.Error("expected SlackURL to be set")
	}
	if cfg.PagerDutyKey == "" {
		t.Error("expected PagerDutyKey to be set")
	}
	if cfg.SplunkURL == "" {
		t.Error("expected SplunkURL to be set")
	}
	if cfg.SplunkSource != "vaultwatch-prod" {
		t.Errorf("unexpected SplunkSource: %s", cfg.SplunkSource)
	}
}

func TestAlertConfig_SplunkDefaults(t *testing.T) {
	cfg := AlertConfig{
		SplunkURL:   "https://splunk.example.com:8088/services/collector/event",
		SplunkToken: "mytoken",
	}
	if cfg.SplunkSource != "" {
		t.Errorf("expected empty SplunkSource by default, got %q", cfg.SplunkSource)
	}
}
