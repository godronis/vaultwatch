package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func main() {
	cfgPath := "vaultwatch.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		log.Fatalf("failed to create vault client: %v", err)
	}

	alerter, err := buildAlerter(cfg.Alerts)
	if err != nil {
		log.Fatalf("failed to build alerter: %v", err)
	}

	if err := vault.Monitor(client, cfg.Vault.Paths, cfg.Vault.WarnBefore, alerter); err != nil {
		log.Fatalf("monitor error: %v", err)
	}
}

func buildAlerter(ac config.AlertConfig) (alert.Alerter, error) {
	var alerters []alert.Alerter

	if ac.Webhook != nil {
		a, err := alert.NewWebhookAlerter(ac.Webhook.URL, ac.Webhook.Headers)
		if err != nil {
			return nil, fmt.Errorf("webhook: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Slack != nil {
		a, err := alert.NewSlackAlerter(ac.Slack.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("slack: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.PagerDuty != nil {
		a, err := alert.NewPagerDutyAlerter(ac.PagerDuty.IntegrationKey)
		if err != nil {
			return nil, fmt.Errorf("pagerduty: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.PagerDutyV2 != nil {
		a, err := alert.NewPagerDutyV2Alerter(ac.PagerDutyV2.RoutingKey)
		if err != nil {
			return nil, fmt.Errorf("pagerduty_v2: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.OpsGenie != nil {
		a, err := alert.NewOpsGenieAlerter(ac.OpsGenie.APIKey)
		if err != nil {
			return nil, fmt.Errorf("opsgenie: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Teams != nil {
		a, err := alert.NewTeamsAlerter(ac.Teams.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("teams: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Datadog != nil {
		a, err := alert.NewDatadogAlerter(ac.Datadog.APIKey)
		if err != nil {
			return nil, fmt.Errorf("datadog: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Telegram != nil {
		a, err := alert.NewTelegramAlerter(ac.Telegram.BotToken, ac.Telegram.ChatID)
		if err != nil {
			return nil, fmt.Errorf("telegram: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.VictorOps != nil {
		a, err := alert.NewVictorOpsAlerter(ac.VictorOps.RESTEndpoint, ac.VictorOps.RoutingKey)
		if err != nil {
			return nil, fmt.Errorf("victorops: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.VictorOpsV2 != nil {
		a, err := alert.NewVictorOpsV2Alerter(ac.VictorOpsV2.RESTEndpoint, ac.VictorOpsV2.RoutingKey)
		if err != nil {
			return nil, fmt.Errorf("victorops_v2: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Discord != nil {
		a, err := alert.NewDiscordAlerter(ac.Discord.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("discord: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.GoogleChat != nil {
		a, err := alert.NewGoogleChatAlerter(ac.GoogleChat.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("googlechat: %w", err)
		}
		alerters = append(alerters, a)
	}

	if ac.Splunk != nil {
		a, err := alert.NewSplunkAlerter(ac.Splunk.URL, ac.Splunk.Token, ac.Splunk.Source, ac.Splunk.SourceType, ac.Splunk.Index)
		if err != nil {
			return nil, fmt.Errorf("splunk: %w", err)
		}
		alerters = append(alerters, a)
	}

	if len(alerters) == 0 {
		return nil, fmt.Errorf("no alerters configured")
	}

	return alert.NewMultiAlerter(alerters...), nil
}
