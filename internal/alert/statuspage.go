package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StatusPageAlerter struct {
	apiKey     string
	pageID     string
	componentID string
	client     *http.Client
}

func NewStatusPageAlerter(apiKey, pageID, componentID string) (*StatusPageAlerter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("statuspage: api key must not be empty")
	}
	if pageID == "" {
		return nil, fmt.Errorf("statuspage: page ID must not be empty")
	}
	if componentID == "" {
		return nil, fmt.Errorf("statuspage: component ID must not be empty")
	}
	return &StatusPageAlerter{
		apiKey:      apiKey,
		pageID:      pageID,
		componentID: componentID,
		client:      &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (s *StatusPageAlerter) Send(path string, expiry time.Time) error {
	url := fmt.Sprintf("https://api.statuspage.io/v1/pages/%s/incidents", s.pageID)

	incident := map[string]interface{}{
		"incident": map[string]interface{}{
			"name": fmt.Sprintf("Vault secret expiring: %s", path),
			"body": fmt.Sprintf(
				"Secret at path '%s' is expiring at %s. Immediate rotation required.",
				path, expiry.Format(time.RFC3339),
			),
			"status":                    "investigating",
			"impact_override":           "minor",
			"component_ids":             []string{s.componentID},
			"deliver_notifications":     true,
		},
	}

	body, err := json.Marshal(incident)
	if err != nil {
		return fmt.Errorf("statuspage: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("statuspage: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "OAuth "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
