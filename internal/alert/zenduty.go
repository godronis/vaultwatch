package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const zendutyAPIURL = "https://www.zenduty.com/api/incidents/"

// ZendutyAlerter sends alerts to Zenduty via their REST API.
type ZendutyAlerter struct {
	apiKey      string
	serviceID   string
	escalationID string
	client      *http.Client
}

// NewZendutyAlerter creates a new ZendutyAlerter.
func NewZendutyAlerter(apiKey, serviceID, escalationID string) (*ZendutyAlerter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("zenduty: api key must not be empty")
	}
	if serviceID == "" {
		return nil, fmt.Errorf("zenduty: service ID must not be empty")
	}
	return &ZendutyAlerter{
		apiKey:       apiKey,
		serviceID:    serviceID,
		escalationID: escalationID,
		client:       &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an alert to Zenduty for the given secret path and expiry.
func (z *ZendutyAlerter) Send(path string, expiresAt time.Time) error {
	payload := map[string]interface{}{
		"title":        fmt.Sprintf("Vault secret expiring: %s", path),
		"summary":      fmt.Sprintf("Secret at path '%s' expires at %s", path, expiresAt.UTC().Format(time.RFC3339)),
		"service":      z.serviceID,
		"urgency":      1,
		"creation_date": time.Now().UTC().Format(time.RFC3339),
	}
	if z.escalationID != "" {
		payload["escalation_policy"] = z.escalationID
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zenduty: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, zendutyAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zenduty: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+z.apiKey)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zenduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
