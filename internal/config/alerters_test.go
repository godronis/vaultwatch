package config

import (
	"testing"
)

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	var ac AlertConfig
	if ac.Slack != nil || ac.PagerDuty != nil || ac.OpsGenie != nil {
		t.Error("expected all alerter configs to be nil by default")
	}
	if ac.VictorOpsV2 != nil {
		t.Error("expected VictorOpsV2 to be nil by default")
	}
}

func TestAlertConfig_TelegramFields(t *testing.T) {
	ac := AlertConfig{
		Telegram: &TelegramConfig{
			BotToken: "bot123",
			ChatID:   "-100123456",
		},
	}
	if ac.Telegram.BotToken != "bot123" {
		t.Errorf("unexpected bot token: %s", ac.Telegram.BotToken)
	}
	if ac.Telegram.ChatID != "-100123456" {
		t.Errorf("unexpected chat id: %s", ac.Telegram.ChatID)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	ac := AlertConfig{
		Slack:     &SlackConfig{WebhookURL: "https://hooks.slack.com/"},
		PagerDuty: &PagerDutyConfig{IntegrationKey: "abc123"},
		Email: &EmailConfig{
			Host:       "smtp.example.com",
			From:       "alerts@example.com",
			Recipients: []string{"ops@example.com"},
		},
	}
	if ac.Slack == nil || ac.PagerDuty == nil || ac.Email == nil {
		t.Error("expected multiple alerters to be set")
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
		t.Errorf("expected empty source, got %s", ac.Splunk.Source)
	}
	if ac.Splunk.Index != "" {
		t.Errorf("expected empty index, got %s", ac.Splunk.Index)
	}
}

func TestAlertConfig_JiraDefaults(t *testing.T) {
	ac := AlertConfig{
		Jira: &JiraConfig{
			URL:     "https://jira.example.com",
			Token:   "token",
			Project: "OPS",
		},
	}
	if ac.Jira.IssueType != "" {
		t.Errorf("expected empty issue type, got %s", ac.Jira.IssueType)
	}
}

func TestAlertConfig_VictorOpsV2Fields(t *testing.T) {
	ac := AlertConfig{
		VictorOpsV2: &VictorOpsV2Config{
			RESTEndpoint: "https://alert.victorops.com/integrations/generic/v2",
			RoutingKey:   "my-routing-key",
		},
	}
	if ac.VictorOpsV2 == nil {
		t.Fatal("expected VictorOpsV2 to be set")
	}
	if ac.VictorOpsV2.RESTEndpoint == "" {
		t.Error("expected non-empty REST endpoint")
	}
	if ac.VictorOpsV2.RoutingKey == "" {
		t.Error("expected non-empty routing key")
	}
}
