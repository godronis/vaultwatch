package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DiscordAlerter sends alerts to a Discord channel via webhook.
type DiscordAlerter struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content  string         `json:"content,omitempty"`
	Embeds   []discordEmbed `json:"embeds,omitempty"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Timestamp   string `json:"timestamp"`
}

// NewDiscordAlerter creates a new DiscordAlerter.
// Returns an error if the webhook URL is empty.
func NewDiscordAlerter(webhookURL string) (*DiscordAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("discord webhook URL must not be empty")
	}
	return &DiscordAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers an alert message to the configured Discord webhook.
func (d *DiscordAlerter) Send(path, message string) error {
	payload := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("VaultWatch: Secret Expiring — %s", path),
				Description: message,
				Color:       15158332, // red
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: failed to marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
