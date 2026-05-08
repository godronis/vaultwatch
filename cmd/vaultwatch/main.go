package main

import (
	"fmt"
	"log"
	"os"

	"github.com/your-org/vaultwatch/internal/alert"
	"github.com/your-org/vaultwatch/internal/config"
	"github.com/your-org/vaultwatch/internal/vault"
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

func buildAlerter(cfg config.AlertConfig) (alert.Alerter, error) {
	var alerters []alert.Alerter

	if cfg.Webhook != nil {
		a, err := alert.NewWebhookAlerter(cfg.Webhook.URL, cfg.Webhook.Headers)
		if err != nil {
			return nil, fmt.Errorf("webhook alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Slack != nil {
		a, err := alert.NewSlackAlerter(cfg.Slack.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("slack alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.PagerDuty != nil {
		a, err := alert.NewPagerDutyAlerter(cfg.PagerDuty.APIKey)
		if err != nil {
			return nil, fmt.Errorf("pagerduty alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.PagerDutyV2 != nil {
		a, err := alert.NewPagerDutyV2Alerter(cfg.PagerDutyV2.IntegrationKey)
		if err != nil {
			return nil, fmt.Errorf("pagerduty v2 alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.OpsGenie != nil {
		a, err := alert.NewOpsGenieAlerter(cfg.OpsGenie.APIKey)
		if err != nil {
			return nil, fmt.Errorf("opsgenie alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Email != nil {
		a, err := alert.NewEmailAlerter(cfg.Email.Host, cfg.Email.Port, cfg.Email.From, cfg.Email.Recipients, cfg.Email.Username, cfg.Email.Password)
		if err != nil {
			return nil, fmt.Errorf("email alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Teams != nil {
		a, err := alert.NewTeamsAlerter(cfg.Teams.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("teams alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Telegram != nil {
		a, err := alert.NewTelegramAlerter(cfg.Telegram.Token, cfg.Telegram.ChatID)
		if err != nil {
			return nil, fmt.Errorf("telegram alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Discord != nil {
		a, err := alert.NewDiscordAlerter(cfg.Discord.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("discord alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Splunk != nil {
		a, err := alert.NewSplunkAlerter(cfg.Splunk.URL, cfg.Splunk.Token, cfg.Splunk.Source, cfg.Splunk.SourceType, cfg.Splunk.Index)
		if err != nil {
			return nil, fmt.Errorf("splunk alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if len(alerters) == 0 {
		return nil, fmt.Errorf("no alerters configured")
	}

	return alert.NewMultiAlerter(alerters...), nil
}
