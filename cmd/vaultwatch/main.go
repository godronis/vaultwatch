// main is the entry point for the vaultwatch CLI tool.
// It loads configuration, initializes alerters, and runs the Vault secret monitor.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/your-org/vaultwatch/internal/alert"
	"github.com/your-org/vaultwatch/internal/config"
	"github.com/your-org/vaultwatch/internal/vault"
)

const defaultConfigPath = "vaultwatch.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "path to configuration file")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", *configPath)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if *verbose {
		log.Printf("loaded config: vault address=%s, warn_before=%s, paths=%d",
			cfg.Vault.Address, cfg.WarnBefore, len(cfg.Paths))
	}

	vaultClient, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		log.Fatalf("failed to create vault client: %v", err)
	}

	alerter, err := buildAlerter(cfg)
	if err != nil {
		log.Fatalf("failed to initialize alerters: %v", err)
	}

	if err := vault.Monitor(vaultClient, alerter, cfg.Paths, cfg.WarnBefore); err != nil {
		log.Fatalf("monitor run failed: %v", err)
	}

	if *verbose {
		log.Println("vaultwatch completed successfully")
	}
}

// buildAlerter constructs a MultiAlerter from the configured alert backends.
// Returns an error if no alerters are configured or any alerter fails to initialize.
func buildAlerter(cfg *config.Config) (alert.Alerter, error) {
	var alerters []alert.Alerter

	if a := cfg.Alerts; a != nil {
		if a.Webhook != nil {
			w, err := alert.NewWebhookAlerter(a.Webhook.URL, a.Webhook.Headers)
			if err != nil {
				return nil, fmt.Errorf("webhook alerter: %w", err)
			}
			alerters = append(alerters, w)
		}

		if a.Slack != nil {
			s, err := alert.NewSlackAlerter(a.Slack.WebhookURL)
			if err != nil {
				return nil, fmt.Errorf("slack alerter: %w", err)
			}
			alerters = append(alerters, s)
		}

		if a.PagerDuty != nil {
			p, err := alert.NewPagerDutyAlerter(a.PagerDuty.IntegrationKey)
			if err != nil {
				return nil, fmt.Errorf("pagerduty alerter: %w", err)
			}
			alerters = append(alerters, p)
		}

		if a.OpsGenie != nil {
			o, err := alert.NewOpsGenieAlerter(a.OpsGenie.APIKey)
			if err != nil {
				return nil, fmt.Errorf("opsgenie alerter: %w", err)
			}
			alerters = append(alerters, o)
		}

		if a.Email != nil {
			e, err := alert.NewEmailAlerter(a.Email.Host, a.Email.Port, a.Email.From, a.Email.To)
			if err != nil {
				return nil, fmt.Errorf("email alerter: %w", err)
			}
			alerters = append(alerters, e)
		}

		if a.Teams != nil {
			t, err := alert.NewTeamsAlerter(a.Teams.WebhookURL)
			if err != nil {
				return nil, fmt.Errorf("teams alerter: %w", err)
			}
			alerters = append(alerters, t)
		}

		if a.Telegram != nil {
			tg, err := alert.NewTelegramAlerter(a.Telegram.BotToken, a.Telegram.ChatID)
			if err != nil {
				return nil, fmt.Errorf("telegram alerter: %w", err)
			}
			alerters = append(alerters, tg)
		}

		if a.Discord != nil {
			d, err := alert.NewDiscordAlerter(a.Discord.WebhookURL)
			if err != nil {
				return nil, fmt.Errorf("discord alerter: %w", err)
			}
			alerters = append(alerters, d)
		}
	}

	if len(alerters) == 0 {
		return nil, fmt.Errorf("no alerters configured: at least one alert backend must be specified")
	}

	return alert.NewMultiAlerter(alerters...), nil
}
