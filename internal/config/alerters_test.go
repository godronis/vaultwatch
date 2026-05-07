package config

import (
	"testing"
)

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	var cfg AlertConfig
	if cfg.Webhook != nil {
		t.Error("expected Webhook to be nil")
	}
	if cfg.Slack != nil {
		t.Error("expected Slack to be nil")
	}
	if cfg.Grafana != nil {
		t.Error("expected Grafana to be nil")
	}
}

func TestAlertConfig_TelegramFields(t *testing.T) {
	cfg := AlertConfig{
		Telegram: &TelegramConfig{
			Token:  "bot-token",
			ChatID: "-100123456",
		},
	}
	if cfg.Telegram.Token != "bot-token" {
		t.Errorf("expected token 'bot-token', got %s", cfg.Telegram.Token)
	}
	if cfg.Telegram.ChatID != "-100123456" {
		t.Errorf("expected chat_id '-100123456', got %s", cfg.Telegram.ChatID)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	cfg := AlertConfig{
		Slack:     &SlackConfig{WebhookURL: "https://hooks.slack.com/xxx"},
		PagerDuty: &PagerDutyConfig{IntegrationKey: "abc123"},
		Grafana:   &GrafanaConfig{WebhookURL: "https://grafana.example.com/webhook"},
	}
	if cfg.Slack == nil {
		t.Error("expected Slack to be set")
	}
	if cfg.PagerDuty == nil {
		t.Error("expected PagerDuty to be set")
	}
	if cfg.Grafana == nil {
		t.Error("expected Grafana to be set")
	}
	if cfg.Webhook != nil {
		t.Error("expected Webhook to remain nil")
	}
}

func TestAlertConfig_SplunkDefaults(t *testing.T) {
	cfg := SplunkConfig{
		URL:   "https://splunk.example.com",
		Token: "splunk-hec-token",
	}
	if cfg.Source != "" {
		t.Errorf("expected empty Source by default, got %s", cfg.Source)
	}
	if cfg.Index != "" {
		t.Errorf("expected empty Index by default, got %s", cfg.Index)
	}
}

func TestAlertConfig_JiraDefaults(t *testing.T) {
	cfg := JiraConfig{
		URL:     "https://jira.example.com",
		Token:   "jira-token",
		Project: "OPS",
	}
	if cfg.IssueType != "" {
		t.Errorf("expected empty IssueType by default, got %s", cfg.IssueType)
	}
}

func TestAlertConfig_GrafanaField(t *testing.T) {
	cfg := AlertConfig{
		Grafana: &GrafanaConfig{
			WebhookURL: "https://grafana.example.com/api/alerts",
		},
	}
	if cfg.Grafana.WebhookURL != "https://grafana.example.com/api/alerts" {
		t.Errorf("unexpected WebhookURL: %s", cfg.Grafana.WebhookURL)
	}
}
