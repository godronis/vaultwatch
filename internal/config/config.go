package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Vault   VaultConfig    `yaml:"vault"`
	Alerts  AlertsConfig   `yaml:"alerts"`
	Monitor MonitorConfig  `yaml:"monitor"`
}

type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Paths   []string `yaml:"paths"`
}

type MonitorConfig struct {
	WarnBefore time.Duration `yaml:"warn_before"`
	Interval   time.Duration `yaml:"interval"`
}

type AlertsConfig struct {
	Webhook   *WebhookConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig  `yaml:"opsgenie,omitempty"`
	Email     *EmailConfig     `yaml:"email,omitempty"`
	Teams     *TeamsConfig     `yaml:"teams,omitempty"`
	Datadog   *DatadogConfig   `yaml:"datadog,omitempty"`
}

type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type PagerDutyConfig struct {
	IntegrationKey string `yaml:"integration_key"`
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

type TeamsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type DatadogConfig struct {
	APIKey string `yaml:"api_key"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse YAML: %w", err)
	}

	if cfg.Monitor.WarnBefore == 0 {
		cfg.Monitor.WarnBefore = 7 * 24 * time.Hour
	}
	if cfg.Monitor.Interval == 0 {
		cfg.Monitor.Interval = 1 * time.Hour
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	var errs []string
	if cfg.Vault.Address == "" {
		errs = append(errs, "vault.address is required")
	}
	if cfg.Vault.Token == "" {
		errs = append(errs, "vault.token is required")
	}
	if len(cfg.Vault.Paths) == 0 {
		errs = append(errs, "vault.paths must contain at least one path")
	}
	if len(errs) > 0 {
		return errors.New("config validation failed: " + joinStrings(errs))
	}
	return nil
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += "; "
		}
		out += s
	}
	return out
}
