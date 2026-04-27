package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackPayload represents a Slack incoming webhook message.
type SlackPayload struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment is a single Slack message attachment.
type SlackAttachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
}

// SlackAlerter sends expiration alerts to a Slack webhook URL.
type SlackAlerter struct {
	WebhookURL string
	Client     *http.Client
}

// NewSlackAlerter creates a new SlackAlerter with the given webhook URL.
func NewSlackAlerter(webhookURL string) (*SlackAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("slack webhook URL must not be empty")
	}
	return &SlackAlerter{
		WebhookURL: webhookURL,
		Client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers a Slack alert for the given secret path and expiry time.
func (s *SlackAlerter) Send(path string, expiresAt time.Time) error {
	payload := SlackPayload{
		Text: ":warning: Vault Secret Expiring Soon",
		Attachments: []SlackAttachment{
			{
				Color: "warning",
				Title: path,
				Text:  fmt.Sprintf("Expires at: %s", expiresAt.UTC().Format(time.RFC3339)),
				Footer: "vaultwatch",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	resp, err := s.Client.Post(s.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
