package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

type DatadogAlerter struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

type datadogPayload struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"`
	Tags      []string `json:"tags,omitempty"`
}

func NewDatadogAlerter(apiKey string) (*DatadogAlerter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("datadog api key must not be empty")
	}
	return &DatadogAlerter{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: datadogEventsURL,
	}, nil
}

func (d *DatadogAlerter) Send(path string, expiresAt time.Time) error {
	payload := datadogPayload{
		Title:     "Vault Secret Expiring Soon",
		Text:      fmt.Sprintf("Secret at path '%s' expires at %s", path, expiresAt.Format(time.RFC3339)),
		AlertType: "warning",
		Tags:      []string{"source:vaultwatch", fmt.Sprintf("secret_path:%s", path)},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("datadog: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
