package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GoogleChatAlerter sends alerts to a Google Chat webhook URL.
type GoogleChatAlerter struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatAlerter creates a new GoogleChatAlerter.
// Returns an error if the webhook URL is empty.
func NewGoogleChatAlerter(webhookURL string) (*GoogleChatAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechat: webhook URL must not be empty")
	}
	return &GoogleChatAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send sends an alert message to the configured Google Chat space.
func (g *GoogleChatAlerter) Send(path string, expiry time.Time) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf(
			"*VaultWatch Alert*\nSecret `%s` is expiring on *%s*.",
			path,
			expiry.UTC().Format(time.RFC3339),
		),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: failed to marshal payload: %w", err)
	}

	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
