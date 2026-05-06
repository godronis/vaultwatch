package config

import (
	"testing"
)

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	var ac AlertConfig
	if ac.Webhook != nil || ac.Slack != nil || ac.PagerDuty != nil {
		t.Error("expected all alerters to be nil by default")
	}
	if ac.CircleCI != nil {
		t.Error("expected CircleCI to be nil by default")
	}
}

func TestAlertConfig_TelegramFields(t *testing.T) {
	ac := AlertConfig{
		Telegram: &TelegramConfig{
			Token:  "bot123",
			ChatID: "-100456",
		},
	}
	if ac.Telegram.Token != "bot123" {
		t.Errorf("unexpected token: %s", ac.Telegram.Token)
	}
	if ac.Telegram.ChatID != "-100456" {
		t.Errorf("unexpected chat_id: %s", ac.Telegram.ChatID)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	ac := AlertConfig{
		Slack:     &SlackConfig{WebhookURL: "https://hooks.slack.com/x"},
		PagerDuty: &PagerDutyConfig{IntegrationKey: "abc"},
		CircleCI:  &CircleCIConfig{Token: "mytoken"},
	}
	if ac.Slack == nil {
		t.Error("expected Slack to be set")
	}
	if ac.PagerDuty == nil {
		t.Error("expected PagerDuty to be set")
	}
	if ac.CircleCI == nil {
		t.Error("expected CircleCI to be set")
	}
	if ac.CircleCI.Token != "mytoken" {
		t.Errorf("unexpected CircleCI token: %s", ac.CircleCI.Token)
	}
}

func TestAlertConfig_SplunkDefaults(t *testing.T) {
	ac := AlertConfig{
		Splunk: &SplunkConfig{
			URL:   "https://splunk.example.com",
			Token: "hec-token",
		},
	}
	if ac.Splunk.Source != "" {
		t.Errorf("expected empty source by default, got %s", ac.Splunk.Source)
	}
}

func TestAlertConfig_JiraDefaults(t *testing.T) {
	ac := AlertConfig{
		Jira: &JiraConfig{
			URL:     "https://jira.example.com",
			Token:   "jira-token",
			Project: "OPS",
		},
	}
	if ac.Jira.IssueType != "" {
		t.Errorf("expected empty issue_type by default, got %s", ac.Jira.IssueType)
	}
}
