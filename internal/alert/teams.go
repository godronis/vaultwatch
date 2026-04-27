package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TeamsAlerter sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsAlerter struct {
	webhookURL string
	client     *http.Client
}

type teamsPayload struct {
	Type       string          `json:"@type"`
	Context    string          `json:"@context"`
	ThemeColor string          `json:"themeColor"`
	Summary    string          `json:"summary"`
	Sections   []teamsSection  `json:"sections"`
}

type teamsSection struct {
	ActivityTitle    string       `json:"activityTitle"`
	ActivitySubtitle string       `json:"activitySubtitle"`
	Facts            []teamsFact  `json:"facts"`
	Markdown         bool         `json:"markdown"`
}

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewTeamsAlerter creates a new TeamsAlerter. Returns an error if the webhook URL is empty.
func NewTeamsAlerter(webhookURL string) (*TeamsAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("teams webhook URL must not be empty")
	}
	return &TeamsAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches a secret expiration alert to Microsoft Teams.
func (t *TeamsAlerter) Send(path string, expiresAt time.Time) error {
	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "FF0000",
		Summary:    fmt.Sprintf("Vault secret expiring: %s", path),
		Sections: []teamsSection{
			{
				ActivityTitle:    "⚠️ Vault Secret Expiration Warning",
				ActivitySubtitle: fmt.Sprintf("Secret at path `%s` is expiring soon", path),
				Facts: []teamsFact{
					{Name: "Path", Value: path},
					{Name: "Expires At", Value: expiresAt.UTC().Format(time.RFC3339)},
					{Name: "Time Remaining", Value: time.Until(expiresAt).Round(time.Second).String()},
				},
				Markdown: true,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: failed to marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
