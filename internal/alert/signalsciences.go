package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SignalSciencesAlerter sends alerts to the Signal Sciences / Fastly NGWAF API.
type SignalSciencesAlerter struct {
	corpName  string
	siteName  string
	apiToken  string
	email     string
	httpClient *http.Client
}

// NewSignalSciencesAlerter creates a new SignalSciencesAlerter.
func NewSignalSciencesAlerter(corpName, siteName, email, apiToken string) (*SignalSciencesAlerter, error) {
	if corpName == "" {
		return nil, fmt.Errorf("signalsciences: corp name is required")
	}
	if siteName == "" {
		return nil, fmt.Errorf("signalsciences: site name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("signalsciences: email is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("signalsciences: API token is required")
	}
	return &SignalSciencesAlerter{
		corpName:   corpName,
		siteName:   siteName,
		email:      email,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts a custom alert event to the Signal Sciences API.
func (a *SignalSciencesAlerter) Send(path string, expiry time.Time) error {
	url := fmt.Sprintf(
		"https://dashboard.signalsciences.net/api/v0/corps/%s/sites/%s/events",
		a.corpName, a.siteName,
	)

	payload := map[string]interface{}{
		"message": fmt.Sprintf("Vault secret expiring: %s at %s", path, expiry.UTC().Format(time.RFC3339)),
		"attachments": []map[string]string{
			{
				"path":   path,
				"expiry": expiry.UTC().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalsciences: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalsciences: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-user", a.email)
	req.Header.Set("x-api-token", a.apiToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("signalsciences: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalsciences: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
