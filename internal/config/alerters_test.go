package config

import (
	"testing"
)

func TestAlertConfig_TelegramFields(t *testing.T) {
	cfg := AlertConfig{
		TelegramToken:  "bot123",
		TelegramChatID: "-1001234567890",
	}
	if cfg.TelegramToken != "bot123" {
		t.Errorf("expected TelegramToken 'bot123', got %q", cfg.TelegramToken)
	}
	if cfg.TelegramChatID != "-1001234567890" {
		t.Errorf("expected TelegramChatID '-1001234567890', got %q", cfg.TelegramChatID)
	}
}

func TestAlertConfig_AllNilByDefault(t *testing.T) {
	cfg := AlertConfig{}
	enabled := cfg.EnabledAlerters()
	if len(enabled) != 0 {
		t.Errorf("expected no enabled alerters, got %v", enabled)
	}
}

func TestAlertConfig_MultipleAlertersSet(t *testing.T) {
	cfg := AlertConfig{
		SlackURL:      "https://hooks.slack.com/test",
		PagerDutyKey:  "pdkey",
		VictorOpsURL:  "https://alert.victorops.com/integrations/generic/1/alert/token/route",
	}
	enabled := cfg.EnabledAlerters()
	if len(enabled) != 3 {
		t.Errorf("expected 3 enabled alerters, got %d: %v", len(enabled), enabled)
	}
	has := func(name string) bool {
		for _, e := range enabled {
			if e == name {
				return true
			}
		}
		return false
	}
	for _, name := range []string{"slack", "pagerduty", "victorops"} {
		if !has(name) {
			t.Errorf("expected %q in enabled alerters", name)
		}
	}
}

func TestAlertConfig_VictorOpsURL(t *testing.T) {
	cfg := AlertConfig{
		VictorOpsURL: "https://alert.victorops.com/integrations/generic/1/alert/token/route",
	}
	enabled := cfg.EnabledAlerters()
	if len(enabled) != 1 || enabled[0] != "victorops" {
		t.Errorf("expected [victorops], got %v", enabled)
	}
}
