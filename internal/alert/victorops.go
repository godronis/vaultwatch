package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// VictorOpsAlerter sends alerts to VictorOps (Splunk On-Call) via REST endpoint.
type VictorOpsAlerter struct {
	webhookURL string
	client     *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

// NewVictorOpsAlerter creates a new VictorOpsAlerter.
// webhookURL is the VictorOps REST endpoint URL including routing key.
func NewVictorOpsAlerter(webhookURL string) (*VictorOpsAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("victorops: webhook URL must not be empty")
	}
	return &VictorOpsAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers an alert to VictorOps for the given secret path and expiry message.
func (v *VictorOpsAlerter) Send(path, message string) error {
	payload := victorOpsPayload{
		MessageType:       "CRITICAL",
		EntityID:          path,
		EntityDisplayName: fmt.Sprintf("Vault Secret Expiring: %s", path),
		StateMessage:      message,
		Timestamp:         time.Now().Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: failed to marshal payload: %w", err)
	}

	resp, err := v.client.Post(v.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
