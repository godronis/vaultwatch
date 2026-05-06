package config

// AlertConfig holds optional configuration blocks for each supported alerter.
// Only the alerters that are non-nil will be activated at runtime.
type AlertConfig struct {
	Webhook   *WebhookConfig   `toml:"webhook"`
	Slack     *SlackConfig     `toml:"slack"`
	PagerDuty *PagerDutyConfig `toml:"pagerduty"`
	OpsGenie  *OpsGenieConfig  `toml:"opsgenie"`
	Email     *EmailConfig     `toml:"email"`
	Teams     *TeamsConfig     `toml:"teams"`
	Datadog   *DatadogConfig   `toml:"datadog"`
	SNS       *SNSConfig       `toml:"sns"`
	Telegram  *TelegramConfig  `toml:"telegram"`
	VictorOps *VictorOpsConfig `toml:"victorops"`
	Discord   *DiscordConfig   `toml:"discord"`
	GoogleChat *GoogleChatConfig `toml:"googlechat"`
	Mattermost *MattermostConfig `toml:"mattermost"`
	Splunk    *SplunkConfig    `toml:"splunk"`
	Zenduty   *ZendutyConfig   `toml:"zenduty"`
	NewRelic  *NewRelicConfig  `toml:"newrelic"`
	Jira      *JiraConfig      `toml:"jira"`
}

type WebhookConfig struct {
	URL     string            `toml:"url"`
	Headers map[string]string `toml:"headers"`
}

type SlackConfig struct {
	WebhookURL string `toml:"webhook_url"`
}

type PagerDutyConfig struct {
	IntegrationKey string `toml:"integration_key"`
}

type OpsGenieConfig struct {
	APIKey string `toml:"api_key"`
}

type EmailConfig struct {
	Host       string   `toml:"host"`
	Port       int      `toml:"port"`
	From       string   `toml:"from"`
	Recipients []string `toml:"recipients"`
}

type TeamsConfig struct {
	WebhookURL string `toml:"webhook_url"`
}

type DatadogConfig struct {
	APIKey string `toml:"api_key"`
}

type SNSConfig struct {
	TopicARN string `toml:"topic_arn"`
}

type TelegramConfig struct {
	BotToken string `toml:"bot_token"`
	ChatID   string `toml:"chat_id"`
}

type VictorOpsConfig struct {
	WebhookURL string `toml:"webhook_url"`
}

type DiscordConfig struct {
	WebhookURL string `toml:"webhook_url"`
}

type GoogleChatConfig struct {
	WebhookURL string `toml:"webhook_url"`
}

type MattermostConfig struct {
	WebhookURL string `toml:"webhook_url"`
	Channel    string `toml:"channel"`
	Username   string `toml:"username"`
}

type SplunkConfig struct {
	URL   string `toml:"url"`
	Token string `toml:"token"`
	Source string `toml:"source"`
}

type ZendutyConfig struct {
	APIKey    string `toml:"api_key"`
	ServiceID string `toml:"service_id"`
}

type NewRelicConfig struct {
	AccountID string `toml:"account_id"`
	APIKey    string `toml:"api_key"`
}

// JiraConfig holds settings for the Jira alerter.
type JiraConfig struct {
	BaseURL   string `toml:"base_url"`
	Token     string `toml:"token"`
	Project   string `toml:"project"`
	IssueType string `toml:"issue_type"`
}
