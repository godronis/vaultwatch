package config

// AlertConfig holds configuration for all supported alerter integrations.
// Only the fields relevant to a chosen alerter need to be populated.
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

	// SNS
	SNSTopicARN string `yaml:"sns_topic_arn"`
	SNSRegion   string `yaml:"sns_region"`

	// Telegram
	TelegramToken  string `yaml:"telegram_token"`
	TelegramChatID string `yaml:"telegram_chat_id"`

	// VictorOps
	VictorOpsURL string `yaml:"victorops_url"`

	// Discord
	DiscordURL string `yaml:"discord_url"`

	// Google Chat
	GoogleChatURL string `yaml:"googlechat_url"`

	// Mattermost
	MattermostURL      string `yaml:"mattermost_url"`
	MattermostChannel  string `yaml:"mattermost_channel"`
	MattermostUsername string `yaml:"mattermost_username"`

	// Splunk HEC
	SplunkURL    string `yaml:"splunk_url"`
	SplunkToken  string `yaml:"splunk_token"`
	SplunkSource string `yaml:"splunk_source"`
}
