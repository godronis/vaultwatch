package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyAlerter sends alerts to PagerDuty via the Events API v2.
type PagerDutyAlerter struct {
	integrationKey string
	client         *http.Client
}

// NewPagerDutyAlerter creates a new PagerDutyAlerter.
func NewPagerDutyAlerter(integrationKey string) (*PagerDutyAlerter, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key must not be empty")
	}
	return &PagerDutyAlerter{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}, nil
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
	Custom    map[string]string `json:"custom_details,omitempty"`
}

// Send dispatches a PagerDuty trigger event for the expiring secret.
func (p *PagerDutyAlerter) Send(path string, expiresAt time.Time) error {
	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:   fmt.Sprintf("Vault secret expiring soon: %s", path),
			Source:    "vaultwatch",
			Severity:  "warning",
			Timestamp: expiresAt.UTC().Format(time.RFC3339),
			Custom: map[string]string{
				"path":       path,
				"expires_at": expiresAt.UTC().Format(time.RFC3339),
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: failed to marshal payload: %w", err)
	}

	resp, err := p.client.Post(pagerDutyEventURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
