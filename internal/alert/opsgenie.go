package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const opsGenieAPIURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieAlerter sends alerts to OpsGenie.
type OpsGenieAlerter struct {
	apiKey string
	client *http.Client
	apiURL string
}

// NewOpsGenieAlerter creates a new OpsGenieAlerter.
// Returns an error if the API key is empty.
func NewOpsGenieAlerter(apiKey string) (*OpsGenieAlerter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsgenie api key must not be empty")
	}
	return &OpsGenieAlerter{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
		apiURL: opsGenieAPIURL,
	}, nil
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Details     map[string]string `json:"details,omitempty"`
}

// Send dispatches an alert to OpsGenie for the given secret path and TTL.
func (o *OpsGenieAlerter) Send(path string, ttl time.Duration) error {
	payload := opsGeniePayload{
		Message:     fmt.Sprintf("Vault secret expiring soon: %s", path),
		Description: fmt.Sprintf("Secret at path '%s' expires in %s.", path, ttl.Round(time.Second)),
		Priority:    "P2",
		Details: map[string]string{
			"path": path,
			"ttl":  ttl.Round(time.Second).String(),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
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
		return fmt.Errorf("opsgenie: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
