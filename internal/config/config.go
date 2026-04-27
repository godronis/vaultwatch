package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Secrets []string      `yaml:"secrets"`
	WarnBefore time.Duration `yaml:"warn_before"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

// AlertsConfig holds all alerting channel configurations.
type AlertsConfig struct {
	Webhook   *WebhookConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig  `yaml:"opsgenie,omitempty"`
	Email     *EmailConfig     `yaml:"email,omitempty"`
	Teams     *TeamsConfig     `yaml:"teams,omitempty"`
}

// TeamsConfig holds Microsoft Teams webhook settings.
type TeamsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type PagerDutyConfig struct {
	RoutingKey string `yaml:"routing_key"`
}

type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
	Username   string   `yaml:"username,omitempty"`
	Password   string   `yaml:"password,omitempty"`
}

const defaultWarnBefore = 72 * time.Hour

// Load reads and validates a config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: invalid YAML: %w", err)
	}

	if cfg.WarnBefore == 0 {
		cfg.WarnBefore = defaultWarnBefore
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	var errs []error
	if cfg.Vault.Address == "" {
		errs = append(errs, errors.New("vault.address is required"))
	}
	if cfg.Vault.Token == "" {
		errs = append(errs, errors.New("vault.token is required"))
	}
	if len(errs) > 0 {
		return fmt.Errorf("config validation failed: %v", errs)
	}
	return nil
}
