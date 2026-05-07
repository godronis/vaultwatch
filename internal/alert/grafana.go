package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GrafanaAlerter sends alerts to a Grafana webhook (e.g. Grafana OnCall or Alertmanager).
type GrafanaAlerter struct {
	webhookURL string
	client     *http.Client
}

func NewGrafanaAlerter(webhookURL string) (*GrafanaAlerter, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("grafana: webhook URL must not be empty")
	}
	return &GrafanaAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (g *GrafanaAlerter) Send(path string, expiresAt time.Time) error {
	payload := map[string]interface{}{
		"title":   fmt.Sprintf("Vault Secret Expiring: %s", path),
		"message": fmt.Sprintf("The secret at '%s' will expire at %s.", path, expiresAt.Format(time.RFC3339)),
		"state":   "alerting",
		"tags": map[string]string{
			"source": "vaultwatch",
			"path":   path,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("grafana: failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, g.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("grafana: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
