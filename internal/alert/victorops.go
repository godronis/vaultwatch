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
func NewVictorOpsAlerter(webhookURL string) (*VictorOpsAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("victorops: webhook URL must not be empty")
	}
	return &VictorOpsAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers an alert payload to VictorOps.
func (v *VictorOpsAlerter) Send(p Payload) error {
	body := victorOpsPayload{
		MessageType:       "CRITICAL",
		EntityID:          p.Path,
		EntityDisplayName: fmt.Sprintf("Vault secret expiring: %s", p.Path),
		StateMessage:      fmt.Sprintf("Secret at path '%s' expires at %s", p.Path, p.ExpiresAt),
		Timestamp:         time.Now().Unix(),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("victorops: failed to marshal payload: %w", err)
	}

	resp, err := v.client.Post(v.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("victorops: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
