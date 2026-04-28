package config

// AlertConfig holds configuration for all supported alerter integrations.
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

// EnabledAlerters returns a list of alerter names that have been configured.
func (a *AlertConfig) EnabledAlerters() []string {
	var enabled []string
	if a.WebhookURL != "" {
		enabled = append(enabled, "webhook")
	}
	if a.SlackURL != "" {
		enabled = append(enabled, "slack")
	}
	if a.PagerDutyKey != "" {
		enabled = append(enabled, "pagerduty")
	}
	if a.OpsGenieKey != "" {
		enabled = append(enabled, "opsgenie")
	}
	if a.EmailHost != "" {
		enabled = append(enabled, "email")
	}
	if a.TeamsURL != "" {
		enabled = append(enabled, "teams")
	}
	if a.DatadogAPIKey != "" {
		enabled = append(enabled, "datadog")
	}
	if a.SNSTopicARN != "" {
		enabled = append(enabled, "sns")
	}
	if a.TelegramToken != "" {
		enabled = append(enabled, "telegram")
	}
	if a.VictorOpsURL != "" {
		enabled = append(enabled, "victorops")
	}
	return enabled
}
