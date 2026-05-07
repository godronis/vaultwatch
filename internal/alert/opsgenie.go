package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const opsGenieAPIURL = "https://api.opsgenie.com/v2/alerts"

type OpsGenieAlerter struct {
	apiKey string
	client *http.Client
}

func NewOpsGenieAlerter(apiKey string) (*OpsGenieAlerter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsgenie: api key must not be empty")
	}
	return &OpsGenieAlerter{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (o *OpsGenieAlerter) Send(path string, expiresAt time.Time) error {
	payload := map[string]interface{}{
		"message":     fmt.Sprintf("Vault secret expiring: %s", path),
		"description": fmt.Sprintf("Secret at path '%s' expires at %s", path, expiresAt.Format(time.RFC3339)),
		"priority":    "P2",
		"tags":        []string{"vaultwatch", "secret-expiry"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, opsGenieAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)
	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
