package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// VictorOpsV2Alerter sends alerts to VictorOps (Splunk On-Call) using the REST endpoint v2.
type VictorOpsV2Alerter struct {
	restEndpoint string
	routingKey   string
	client       *http.Client
}

// NewVictorOpsV2Alerter creates a new VictorOpsV2Alerter.
func NewVictorOpsV2Alerter(restEndpoint, routingKey string) (*VictorOpsV2Alerter, error) {
	if restEndpoint == "" {
		return nil, fmt.Errorf("victorops v2: rest endpoint is required")
	}
	if routingKey == "" {
		return nil, fmt.Errorf("victorops v2: routing key is required")
	}
	return &VictorOpsV2Alerter{
		restEndpoint: restEndpoint,
		routingKey:   routingKey,
		client:       &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send sends an alert to VictorOps v2 REST API.
func (a *VictorOpsV2Alerter) Send(path string, expiry time.Time) error {
	payload := map[string]interface{}{
		"message_type":        "CRITICAL",
		"entity_id":          fmt.Sprintf("vaultwatch/%s", path),
		"entity_display_name": fmt.Sprintf("Vault secret expiring: %s", path),
		"state_message":      fmt.Sprintf("Secret at path '%s' expires at %s", path, expiry.Format(time.RFC3339)),
		"state_start_time":   time.Now().Unix(),
		"monitoring_tool":    "vaultwatch",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops v2: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", a.restEndpoint, a.routingKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops v2: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("victorops v2: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops v2: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
