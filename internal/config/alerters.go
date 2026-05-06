package config

// AlertConfig holds configuration for all supported alerter backends.
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

	// Teams
	TeamsURL string `yaml:"teams_url"`

	// Datadog
	DatadogKey  string   `yaml:"datadog_key"`
	DatadogTags []string `yaml:"datadog_tags"`

	// SNS
	SNSTopicARN string `yaml:"sns_topic_arn"`

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

	// Splunk
	SplunkURL   string `yaml:"splunk_url"`
	SplunkToken string `yaml:"splunk_token"`
	SplunkSource string `yaml:"splunk_source"`
	SplunkIndex  string `yaml:"splunk_index"`

	// Zenduty
	ZendutyAPIKey      string `yaml:"zenduty_api_key"`
	ZendutyServiceID   string `yaml:"zenduty_service_id"`
	ZendutyEscalationID string `yaml:"zenduty_escalation_id"`
}
