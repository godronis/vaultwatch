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
	if cfg.PagerDutyV2 != nil {
		t.Error("expected PagerDutyV2 to be nil")
	}
}

func TestAlertConfig_TelegramFields(t *testing.T) {
	cfg := AlertConfig{
		Telegram: &TelegramConfig{
			Token:  "bot-token",
			ChatID: "12345",
		},
	}
	if cfg.Telegram.Token != "bot-token" {
		t.Errorf("expected token 'bot-token', got %s", cfg.Telegram.Token)
	}
	if cfg.Telegram.ChatID != "12345" {
		t.Errorf("expected chat_id '12345', got %s", cfg.Telegram.ChatID)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	cfg := AlertConfig{
		Slack:       &SlackConfig{WebhookURL: "https://hooks.slack.com/test"},
		PagerDutyV2: &PagerDutyV2Config{IntegrationKey: "key-123"},
		Datadog:     &DatadogConfig{APIKey: "dd-key"},
	}
	if cfg.Slack == nil || cfg.PagerDutyV2 == nil || cfg.Datadog == nil {
		t.Error("expected all three alerters to be set")
	}
}

func TestAlertConfig_SplunkDefaults(t *testing.T) {
	cfg := SplunkConfig{
		URL:   "https://splunk.example.com",
		Token: "splunk-token",
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

func TestAlertConfig_PagerDutyV2Fields(t *testing.T) {
	cfg := AlertConfig{
		PagerDutyV2: &PagerDutyV2Config{
			IntegrationKey: "routing-key-abc",
		},
	}
	if cfg.PagerDutyV2.IntegrationKey != "routing-key-abc" {
		t.Errorf("expected integration key 'routing-key-abc', got %s", cfg.PagerDutyV2.IntegrationKey)
	}
}
