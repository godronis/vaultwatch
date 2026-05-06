package config

// AlertConfig holds configuration for all supported alerter integrations.
type AlertConfig struct {
	// Webhook
	WebhookURL     string            `yaml:"webhook_url"`
	WebhookHeaders map[string]string `yaml:"webhook_headers"`

	// Slack
	SlackWebhookURL string `yaml:"slack_webhook_url"`

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
	TeamsWebhookURL string `yaml:"teams_webhook_url"`

	// Datadog
	DatadogAPIKey string `yaml:"datadog_api_key"`

	// SNS
	SNSTopicARN string `yaml:"sns_topic_arn"`

	// Telegram
	TelegramToken  string `yaml:"telegram_token"`
	TelegramChatID string `yaml:"telegram_chat_id"`

	// VictorOps
	VictorOpsURL string `yaml:"victorops_url"`

	// Discord
	DiscordWebhookURL string `yaml:"discord_webhook_url"`

	// Google Chat
	GoogleChatWebhookURL string `yaml:"googlechat_webhook_url"`

	// Mattermost
	MattermostWebhookURL string `yaml:"mattermost_webhook_url"`
	MattermostChannel    string `yaml:"mattermost_channel"`
	MattermostUsername   string `yaml:"mattermost_username"`

	// Splunk
	SplunkURL   string `yaml:"splunk_url"`
	SplunkToken string `yaml:"splunk_token"`
	SplunkSource string `yaml:"splunk_source"`

	// Zenduty
	ZendutyAPIKey    string `yaml:"zenduty_api_key"`
	ZendutyServiceID string `yaml:"zenduty_service_id"`

	// New Relic
	NewRelicAccountID string `yaml:"newrelic_account_id"`
	NewRelicAPIKey    string `yaml:"newrelic_api_key"`
}
