package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const newRelicEventsURL = "https://insights-collector.newrelic.com/v1/accounts/%s/events"

type newRelicAlerter struct {
	accountID string
	apiKey    string
	client    *http.Client
	endpoint  string
}

// NewNewRelicAlerter creates an alerter that sends custom events to New Relic Insights.
func NewNewRelicAlerter(accountID, apiKey string) (*newRelicAlerter, error) {
	if accountID == "" {
		return nil, fmt.Errorf("new relic account ID must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("new relic API key must not be empty")
	}
	return &newRelicAlerter{
		accountID: accountID,
		apiKey:    apiKey,
		client:    &http.Client{Timeout: 10 * time.Second},
		endpoint:  fmt.Sprintf(newRelicEventsURL, accountID),
	}, nil
}

func (a *newRelicAlerter) Send(path string, expiresAt time.Time) error {
	payload := map[string]interface{}{
		"eventType":  "VaultSecretExpiring",
		"secretPath": path,
		"expiresAt":  expiresAt.UTC().Format(time.RFC3339),
		"message":    fmt.Sprintf("Vault secret at %s expires at %s", path, expiresAt.UTC().Format(time.RFC3339)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("newrelic: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, a.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", a.apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
