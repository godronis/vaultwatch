package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Monitor MonitorConfig `yaml:"monitor"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Paths   []string `yaml:"paths"`
}

// MonitorConfig controls polling behaviour.
type MonitorConfig struct {
	Interval   time.Duration `yaml:"interval"`
	WarnBefore time.Duration `yaml:"warn_before"`
}

// AlertsConfig lists all optional alerting back-ends.
type AlertsConfig struct {
	Webhook   *WebhookAlertConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackAlertConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyAlertConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieAlertConfig  `yaml:"opsgenie,omitempty"`
	Email     *EmailAlertConfig     `yaml:"email,omitempty"`
	Teams     *TeamsAlertConfig     `yaml:"teams,omitempty"`
	Datadog   *DatadogAlertConfig   `yaml:"datadog,omitempty"`
	SNS       *SNSAlertConfig       `yaml:"sns,omitempty"`
}

// SNSAlertConfig holds AWS SNS alerting configuration.
type SNSAlertConfig struct {
	TopicARN string `yaml:"topic_arn"`
}

type WebhookAlertConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}
type SlackAlertConfig struct{ URL string `yaml:"url"` }
type PagerDutyAlertConfig struct{ IntegrationKey string `yaml:"integration_key"` }
type OpsGenieAlertConfig struct{ APIKey string `yaml:"api_key"` }
type EmailAlertConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}
type TeamsAlertConfig struct{ URL string `yaml:"url"` }
type DatadogAlertConfig struct{ APIKey string `yaml:"api_key"` }

// Load reads and validates a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: invalid YAML: %w", err)
	}

	if cfg.Monitor.WarnBefore == 0 {
		cfg.Monitor.WarnBefore = 72 * time.Hour
	}
	if cfg.Monitor.Interval == 0 {
		cfg.Monitor.Interval = 5 * time.Minute
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
		return errors.New("config validation failed: " + joinStrings(errs, "; "))
	}
	return nil
}

func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}
