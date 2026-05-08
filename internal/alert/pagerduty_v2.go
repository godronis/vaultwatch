package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyV2EventsURL = "https://events.pagerduty.com/v2/enqueue"

type PagerDutyV2Alerter struct {
	integrationKey string
	client         *http.Client
}

func NewPagerDutyV2Alerter(integrationKey string) (*PagerDutyV2Alerter, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty v2: integration key must not be empty")
	}
	return &PagerDutyV2Alerter{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (a *PagerDutyV2Alerter) Send(path string, expiry time.Time) error {
	payload := map[string]interface{}{
		"routing_key":  a.integrationKey,
		"event_action": "trigger",
		"payload": map[string]interface{}{
			"summary":   fmt.Sprintf("Vault secret expiring: %s", path),
			"severity":  "warning",
			"source":    "vaultwatch",
			"timestamp": expiry.UTC().Format(time.RFC3339),
			"custom_details": map[string]string{
				"secret_path": path,
				"expires_at":  expiry.UTC().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty v2: marshal payload: %w", err)
	}

	resp, err := a.client.Post(pagerDutyV2EventsURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagerduty v2: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty v2: unexpected status %d", resp.StatusCode)
	}
	return nil
}
