package config

import (
	"testing"
)

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	var ac AlertConfig
	if ac.Webhook != nil || ac.Slack != nil || ac.PagerDuty != nil {
		t.Error("expected all alerter configs to be nil by default")
	}
	if ac.Jira != nil {
		t.Error("expected Jira config to be nil by default")
	}
}

func TestAlertConfig_TelegramFields(t *testing.T) {
	ac := AlertConfig{
		Telegram: &TelegramConfig{
			BotToken: "abc123",
			ChatID:   "-100999",
		},
	}
	if ac.Telegram.BotToken != "abc123" {
		t.Errorf("unexpected bot token: %q", ac.Telegram.BotToken)
	}
	if ac.Telegram.ChatID != "-100999" {
		t.Errorf("unexpected chat ID: %q", ac.Telegram.ChatID)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	ac := AlertConfig{
		Slack:     &SlackConfig{WebhookURL: "https://hooks.slack.com/test"},
		PagerDuty: &PagerDutyConfig{IntegrationKey: "key123"},
		Jira:      &JiraConfig{BaseURL: "https://jira.example.com", Token: "t", Project: "OPS"},
	}
	if ac.Slack == nil || ac.PagerDuty == nil || ac.Jira == nil {
		t.Error("expected Slack, PagerDuty, and Jira configs to be set")
	}
}

func TestAlertConfig_SplunkDefaults(t *testing.T) {
	ac := AlertConfig{
		Splunk: &SplunkConfig{
			URL:   "https://splunk.example.com",
			Token: "splunk-token",
		},
	}
	if ac.Splunk.Source != "" {
		t.Errorf("expected empty source by default, got %q", ac.Splunk.Source)
	}
}

func TestAlertConfig_JiraDefaults(t *testing.T) {
	ac := AlertConfig{
		Jira: &JiraConfig{
			BaseURL: "https://jira.example.com",
			Token:   "mytoken",
			Project: "SEC",
		},
	}
	if ac.Jira.IssueType != "" {
		t.Errorf("expected empty IssueType by default, got %q", ac.Jira.IssueType)
	}
	if ac.Jira.Project != "SEC" {
		t.Errorf("expected project 'SEC', got %q", ac.Jira.Project)
	}
}
