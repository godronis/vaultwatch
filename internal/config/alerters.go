package config

// WebhookConfig holds configuration for a generic webhook alerter.
type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// SlackConfig holds configuration for the Slack alerter.
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// PagerDutyConfig holds configuration for the PagerDuty alerter.
type PagerDutyConfig struct {
	IntegrationKey string `yaml:"integration_key"`
}

// OpsGenieConfig holds configuration for the OpsGenie alerter.
type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

// EmailConfig holds configuration for the email alerter.
type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port,omitempty"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// TeamsConfig holds configuration for the Microsoft Teams alerter.
type TeamsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// DatadogConfig holds configuration for the Datadog alerter.
type DatadogConfig struct {
	APIKey string `yaml:"api_key"`
}

// SNSConfig holds configuration for the AWS SNS alerter.
type SNSConfig struct {
	TopicARN string `yaml:"topic_arn"`
}

// TelegramConfig holds configuration for the Telegram alerter.
type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

// VictorOpsConfig holds configuration for the VictorOps alerter.
type VictorOpsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// DiscordConfig holds configuration for the Discord alerter.
type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// AlertConfig aggregates all supported alerter configurations.
type AlertConfig struct {
	Webhook   *WebhookConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig  `yaml:"opsgenie,omitempty"`
	Email     *EmailConfig     `yaml:"email,omitempty"`
	Teams     *TeamsConfig     `yaml:"teams,omitempty"`
	Datadog   *DatadogConfig   `yaml:"datadog,omitempty"`
	SNS       *SNSConfig       `yaml:"sns,omitempty"`
	Telegram  *TelegramConfig  `yaml:"telegram,omitempty"`
	VictorOps *VictorOpsConfig `yaml:"victorops,omitempty"`
	Discord   *DiscordConfig   `yaml:"discord,omitempty"`
}
