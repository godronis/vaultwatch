package config

// AlertConfig holds optional configuration for each supported alerter.
// Only fields for enabled alerters need to be populated.
type AlertConfig struct {
	Webhook       *WebhookConfig       `yaml:"webhook,omitempty"`
	Slack         *SlackConfig         `yaml:"slack,omitempty"`
	PagerDuty     *PagerDutyConfig     `yaml:"pagerduty,omitempty"`
	PagerDutyV2   *PagerDutyV2Config   `yaml:"pagerduty_v2,omitempty"`
	OpsGenie      *OpsGenieConfig      `yaml:"opsgenie,omitempty"`
	Email         *EmailConfig         `yaml:"email,omitempty"`
	Teams         *TeamsConfig         `yaml:"teams,omitempty"`
	Datadog       *DatadogConfig       `yaml:"datadog,omitempty"`
	SNS           *SNSConfig           `yaml:"sns,omitempty"`
	Telegram      *TelegramConfig      `yaml:"telegram,omitempty"`
	VictorOps     *VictorOpsConfig     `yaml:"victorops,omitempty"`
	Discord       *DiscordConfig       `yaml:"discord,omitempty"`
	GoogleChat    *GoogleChatConfig    `yaml:"googlechat,omitempty"`
	Mattermost    *MattermostConfig    `yaml:"mattermost,omitempty"`
	Splunk        *SplunkConfig        `yaml:"splunk,omitempty"`
	Zenduty       *ZendutyConfig       `yaml:"zenduty,omitempty"`
	NewRelic      *NewRelicConfig      `yaml:"newrelic,omitempty"`
	Jira          *JiraConfig          `yaml:"jira,omitempty"`
	SignalSciences *SignalSciencesConfig `yaml:"signalsciences,omitempty"`
	CircleCI      *CircleCIConfig      `yaml:"circleci,omitempty"`
	Grafana       *GrafanaConfig       `yaml:"grafana,omitempty"`
	RocketChat    *RocketChatConfig    `yaml:"rocketchat,omitempty"`
	StatusPage    *StatusPageConfig    `yaml:"statuspage,omitempty"`
}

type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type PagerDutyConfig struct {
	APIKey string `yaml:"api_key"`
}

type PagerDutyV2Config struct {
	IntegrationKey string `yaml:"integration_key"`
}

type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
	Username   string   `yaml:"username,omitempty"`
	Password   string   `yaml:"password,omitempty"`
}

type TeamsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type DatadogConfig struct {
	APIKey string `yaml:"api_key"`
}

type SNSConfig struct {
	TopicARN string `yaml:"topic_arn"`
}

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID string `yaml:"chat_id"`
}

type VictorOpsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type GoogleChatConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type MattermostConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel,omitempty"`
	Username   string `yaml:"username,omitempty"`
}

type SplunkConfig struct {
	URL        string `yaml:"url"`
	Token      string `yaml:"token"`
	Source     string `yaml:"source,omitempty"`
	SourceType string `yaml:"sourcetype,omitempty"`
	Index      string `yaml:"index,omitempty"`
}

type ZendutyConfig struct {
	APIKey    string `yaml:"api_key"`
	ServiceID string `yaml:"service_id"`
}

type NewRelicConfig struct {
	AccountID string `yaml:"account_id"`
	APIKey    string `yaml:"api_key"`
}

type JiraConfig struct {
	URL       string `yaml:"url"`
	Token     string `yaml:"token"`
	Project   string `yaml:"project"`
	IssueType string `yaml:"issue_type,omitempty"`
}

type SignalSciencesConfig struct {
	CorpName string `yaml:"corp_name"`
	SiteName string `yaml:"site_name"`
	Email    string `yaml:"email"`
	Token    string `yaml:"token"`
}

type CircleCIConfig struct {
	Token string `yaml:"token"`
}

type GrafanaConfig struct {
	URL string `yaml:"url"`
}

type RocketChatConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type StatusPageConfig struct {
	APIKey      string `yaml:"api_key"`
	PageID      string `yaml:"page_id"`
	ComponentID string `yaml:"component_id"`
}
