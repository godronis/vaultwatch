package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SplunkAlerter sends alerts to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkAlerter struct {
	url   string
	token string
	source string
	client *http.Client
}

type splunkEvent struct {
	Time   int64                  `json:"time"`
	Source string                 `json:"source"`
	Event  map[string]interface{} `json:"event"`
}

// NewSplunkAlerter creates a new SplunkAlerter.
// url should be the HEC endpoint, e.g. https://splunk.example.com:8088/services/collector/event
func NewSplunkAlerter(url, token, source string) (*SplunkAlerter, error) {
	if url == "" {
		return nil, fmt.Errorf("splunk: HEC url is required")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: HEC token is required")
	}
	if source == "" {
		source = "vaultwatch"
	}
	return &SplunkAlerter{
		url:    url,
		token:  token,
		source: source,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers an alert event to Splunk HEC.
func (s *SplunkAlerter) Send(path string, expiresAt time.Time) error {
	payload := splunkEvent{
		Time:   time.Now().Unix(),
		Source: s.source,
		Event: map[string]interface{}{
			"alert":      "vault_secret_expiring",
			"path":       path,
			"expires_at": expiresAt.UTC().Format(time.RFC3339),
			"message":    fmt.Sprintf("Vault secret at %s expires at %s", path, expiresAt.UTC().Format(time.RFC3339)),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
