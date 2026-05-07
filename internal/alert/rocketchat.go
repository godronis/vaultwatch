package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RocketChatAlerter sends alert notifications to a Rocket.Chat webhook.
type RocketChatAlerter struct {
	webhookURL string
	client     *http.Client
}

// NewRocketChatAlerter creates a new RocketChatAlerter.
// Returns an error if the webhook URL is empty.
func NewRocketChatAlerter(webhookURL string) (*RocketChatAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("rocketchat: webhook URL must not be empty")
	}
	return &RocketChatAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts an alert message to the configured Rocket.Chat webhook.
func (r *RocketChatAlerter) Send(path string, expiry time.Time) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf(
			":warning: *VaultWatch Alert*\nSecret `%s` expires at *%s*.",
			path,
			expiry.UTC().Format(time.RFC3339),
		),
		"attachments": []map[string]interface{}{
			{
				"title": "Secret Expiration Warning",
				"text":  fmt.Sprintf("Path: %s\nExpires: %s", path, expiry.UTC().Format(time.RFC3339)),
				"color": "#FF0000",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: failed to marshal payload: %w", err)
	}

	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
