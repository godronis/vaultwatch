package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level VaultWatch configuration.
type Config struct {
	Vault    VaultConfig    `yaml:"vault"`
	Alerts   AlertsConfig   `yaml:"alerts"`
	Schedule string         `yaml:"schedule"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string        `yaml:"address"`
	Token     string        `yaml:"token"`
	Namespace string        `yaml:"namespace"`
	Paths     []string      `yaml:"paths"`
	WarnBefore time.Duration `yaml:"warn_before"`
}

// AlertsConfig holds webhook alert destinations.
type AlertsConfig struct {
	Webhooks []WebhookConfig `yaml:"webhooks"`
}

// WebhookConfig defines a single webhook endpoint.
type WebhookConfig struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token is required")
	}
	if len(c.Vault.Paths) == 0 {
		return fmt.Errorf("vault.paths must contain at least one path")
	}
	if c.Vault.WarnBefore == 0 {
		c.Vault.WarnBefore = 7 * 24 * time.Hour // default: 7 days
	}
	for i, wh := range c.Alerts.Webhooks {
		if wh.URL == "" {
			return fmt.Errorf("alerts.webhooks[%d].url is required", i)
		}
		if wh.Timeout == 0 {
			c.Alerts.Webhooks[i].Timeout = 10 * time.Second
		}
	}
	return nil
}
