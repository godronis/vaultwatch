package config

// AlertConfig holds optional configuration for each supported alerter.
// Only fields relevant to the configured alerter type need to be set.
type AlertConfig struct {
	// Webhook
	WebhookURL     string            `yaml:"webhook_url"`
	WebhookHeaders map[string]string `yaml:"webhook_headers"`

	// Slack
	SlackURL string `yaml:"slack_url"`

	// PagerDuty
	PagerDutyKey string `yaml:"pagerduty_key"`

	// OpsGenie
	OpsGenieKey string `yaml:"opsgenie_key"`

	// Email
	EmailHost       string   `yaml:"email_host"`
	EmailPort       int      `yaml:"email_port"`
	EmailFrom       string   `yaml:"email_from"`
	EmailTo         []string `yaml:"email_to"`
	EmailUsername   string   `yaml:"email_username"`
	EmailPassword   string   `yaml:"email_password"`

	// Microsoft Teams
	TeamsURL string `yaml:"teams_url"`

	// Datadog
	DatadogAPIKey string   `yaml:"datadog_api_key"`
	DatadogTags   []string `yaml:"datadog_tags"`

	// AWS SNS
	SNSTopicARN string `yaml:"sns_topic_arn"`
	SNSRegion   string `yaml:"sns_region"`

	// Telegram
	TelegramToken  string `yaml:"telegram_token"`
	TelegramChatID string `yaml:"telegram_chat_id"`

	// VictorOps
	VictorOpsURL string `yaml:"victorops_url"`
}
