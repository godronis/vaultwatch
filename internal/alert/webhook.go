// Package alert provides webhook-based alerting functionality for VaultWatch.
// It supports sending secret expiration notifications to configurable HTTP endpoints.
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the JSON body sent to a webhook endpoint when a secret
// is approaching expiration or has already expired.
type Payload struct {
	SecretPath string    `json:"secret_path"`
	ExpireTime time.Time `json:"expire_time"`
	TTL        int       `json:"ttl_seconds"`
	Severity   string    `json:"severity"` // "warning" or "critical"
	Message    string    `json:"message"`
}

// WebhookAlerter sends alert payloads to one or more webhook URLs.
type WebhookAlerter struct {
	urls       []string
	client     *http.Client
	timeoutSec int
}

// NewWebhookAlerter creates a new WebhookAlerter with the given target URLs.
// A default HTTP timeout of 10 seconds is applied unless overridden via options.
func NewWebhookAlerter(urls []string) *WebhookAlerter {
	return &WebhookAlerter{
		urls: urls,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		timeoutSec: 10,
	}
}

// Send dispatches the alert payload to all configured webhook URLs.
// It returns a combined error if any deliveries fail, but continues
// attempting all endpoints regardless of individual failures.
func (w *WebhookAlerter) Send(payload Payload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alert: failed to marshal payload: %w", err)
	}

	var errs []error
	for _, url := range w.urls {
		if err := w.post(url, body); err != nil {
			errs = append(errs, fmt.Errorf("alert: webhook %s failed: %w", url, err))
		}
	}

	if len(errs) > 0 {
		return joinErrors(errs)
	}
	return nil
}

// post sends the JSON body to a single webhook URL via HTTP POST.
func (w *WebhookAlerter) post(url string, body []byte) error {
	resp, err := w.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// BuildPayload constructs an alert Payload for a given secret path and expiry time.
// The severity is set to "critical" if the secret expires within 1 hour, otherwise "warning".
func BuildPayload(secretPath string, expireTime time.Time, ttlSeconds int) Payload {
	severity := "warning"
	if time.Until(expireTime) <= time.Hour {
		severity = "critical"
	}

	message := fmt.Sprintf(
		"Secret '%s' expires in %s (TTL: %ds)",
		secretPath,
		time.Until(expireTime).Round(time.Second),
		ttlSeconds,
	)

	return Payload{
		SecretPath: secretPath,
		ExpireTime: expireTime,
		TTL:        ttlSeconds,
		Severity:   severity,
		Message:    message,
	}
}

// joinErrors combines multiple errors into a single descriptive error.
func joinErrors(errs []error) error {
	msg := fmt.Sprintf("%d webhook(s) failed:", len(errs))
	for _, e := range errs {
		msg += "\n  - " + e.Error()
	}
	return fmt.Errorf("%s", msg)
}
