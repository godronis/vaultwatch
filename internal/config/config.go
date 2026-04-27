package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	WarnBefore time.Duration `yaml:"warn_before"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Paths   []string `yaml:"paths"`
}

// AlertsConfig holds all alerting backend configurations.
type AlertsConfig struct {
	Webhook   *WebhookConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig  `yaml:"opsgenie,omitempty"`
}

// WebhookConfig holds generic webhook settings.
type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// SlackConfig holds Slack webhook settings.
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// PagerDutyConfig holds PagerDuty integration settings.
type PagerDutyConfig struct {
	IntegrationKey string `yaml:"integration_key"`
}

// OpsGenieConfig holds OpsGenie integration settings.
type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

const defaultWarnBefore = 72 * time.Hour

// Load reads and validates the configuration from the given file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse yaml: %w", err)
	}

	if cfg.Vault.Address == "" {
		return nil, fmt.Errorf("config: vault.address is required")
	}
	if cfg.Vault.Token == "" {
		return nil, fmt.Errorf("config: vault.token is required")
	}
	if cfg.WarnBefore == 0 {
		cfg.WarnBefore = defaultWarnBefore
	}

	return &cfg, nil
}
