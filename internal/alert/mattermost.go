package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MattermostAlerter sends alerts to a Mattermost incoming webhook.
type MattermostAlerter struct {
	webhookURL string
	channel    string
	username   string
	client     *http.Client
}

type mattermostPayload struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
}

// NewMattermostAlerter creates a new MattermostAlerter.
// webhookURL is required; channel and username are optional overrides.
func NewMattermostAlerter(webhookURL, channel, username string) (*MattermostAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("mattermost webhook URL must not be empty")
	}
	return &MattermostAlerter{
		webhookURL: webhookURL,
		channel:    channel,
		username:   username,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts an alert message to the configured Mattermost webhook.
func (m *MattermostAlerter) Send(path string, expiresAt time.Time) error {
	text := fmt.Sprintf(
		":warning: **VaultWatch Alert** — secret `%s` expires at `%s`",
		path, expiresAt.UTC().Format(time.RFC3339),
	)

	payload := mattermostPayload{
		Text:     text,
		Channel:  m.channel,
		Username: m.username,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: failed to marshal payload: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
