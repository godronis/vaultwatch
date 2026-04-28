package config

import (
	"testing"
)

func TestAlertConfig_TelegramFields(t *testing.T) {
	cfg := AlertConfig{
		Telegram: &TelegramConfig{
			BotToken: "bot123:ABC",
			ChatID:   "-1001234567890",
		},
	}

	if cfg.Telegram == nil {
		t.Fatal("expected Telegram config to be set")
	}
	if cfg.Telegram.BotToken != "bot123:ABC" {
		t.Errorf("unexpected BotToken: %s", cfg.Telegram.BotToken)
	}
	if cfg.Telegram.ChatID != "-1001234567890" {
		t.Errorf("unexpected ChatID: %s", cfg.Telegram.ChatID)
	}
}

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	cfg := AlertConfig{}

	if cfg.Webhook != nil {
		t.Error("expected Webhook to be nil")
	}
	if cfg.Slack != nil {
		t.Error("expected Slack to be nil")
	}
	if cfg.PagerDuty != nil {
		t.Error("expected PagerDuty to be nil")
	}
	if cfg.Telegram != nil {
		t.Error("expected Telegram to be nil")
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	cfg := AlertConfig{
		Slack: &SlackAlertConfig{WebhookURL: "https://hooks.slack.com/test"},
		Telegram: &TelegramConfig{BotToken: "tok", ChatID: "42"},
		Email: &EmailAlertConfig{
			Host:       "smtp.example.com",
			Port:       587,
			From:       "alerts@example.com",
			Recipients: []string{"admin@example.com"},
		},
	}

	if cfg.Slack == nil || cfg.Telegram == nil || cfg.Email == nil {
		t.Error("expected all three alerters to be configured")
	}
	if cfg.Email.Port != 587 {
		t.Errorf("unexpected email port: %d", cfg.Email.Port)
	}
}
