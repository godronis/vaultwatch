package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const circleCIAPIURL = "https://circleci.com/api/v2/insights/events"

// CircleCIAlerter sends Vault secret expiry notifications as CircleCI pipeline events.
type CircleCIAlerter struct {
	token  string
	client *http.Client
	apiURL string
}

func NewCircleCIAlerter(token string) (*CircleCIAlerter, error) {
	if token == "" {
		return nil, fmt.Errorf("circleci: token must not be empty")
	}
	return &CircleCIAlerter{
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
		apiURL: circleCIAPIURL,
	}, nil
}

func (c *CircleCIAlerter) Send(path string, expiresAt time.Time) error {
	payload := map[string]interface{}{
		"name":       "vault-secret-expiry",
		"attributes": map[string]string{
			"secret_path": path,
			"expires_at":  expiresAt.Format(time.RFC3339),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("circleci: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("circleci: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("circleci: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("circleci: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
