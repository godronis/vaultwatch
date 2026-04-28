package config

// TelegramConfig holds configuration for the Telegram alerter.
type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

// AlertConfig aggregates all supported alerter configurations.
// Each field is optional; only configured alerters will be activated.
type AlertConfig struct {
	Webhook   *WebhookAlertConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackAlertConfig     `yaml:"slack,omitempty"`
	PagerDuty *PagerDutyAlertConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieAlertConfig  `yaml:"opsgenie,omitempty"`
	Email     *EmailAlertConfig     `yaml:"email,omitempty"`
	Teams     *TeamsAlertConfig     `yaml:"teams,omitempty"`
	Datadog   *DatadogAlertConfig   `yaml:"datadog,omitempty"`
	SNS       *SNSAlertConfig       `yaml:"sns,omitempty"`
	Telegram  *TelegramConfig       `yaml:"telegram,omitempty"`
}

// WebhookAlertConfig holds generic webhook alerter settings.
type WebhookAlertConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// SlackAlertConfig holds Slack webhook settings.
type SlackAlertConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// PagerDutyAlertConfig holds PagerDuty integration settings.
type PagerDutyAlertConfig struct {
	IntegrationKey string `yaml:"integration_key"`
}

// OpsGenieAlertConfig holds OpsGenie API settings.
type OpsGenieAlertConfig struct {
	APIKey string `yaml:"api_key"`
}

// EmailAlertConfig holds SMTP email alerter settings.
type EmailAlertConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
	Username   string   `yaml:"username,omitempty"`
	Password   string   `yaml:"password,omitempty"`
}

// TeamsAlertConfig holds Microsoft Teams webhook settings.
type TeamsAlertConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// DatadogAlertConfig holds Datadog API settings.
type DatadogAlertConfig struct {
	APIKey string `yaml:"api_key"`
}

// SNSAlertConfig holds AWS SNS alerter settings.
type SNSAlertConfig struct {
	TopicARN string `yaml:"topic_arn"`
	Region   string `yaml:"region"`
}
