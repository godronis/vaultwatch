package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// JiraAlerter creates Jira issues for expiring secrets.
type JiraAlerter struct {
	baseURL   string
	token     string
	project   string
	issueType string
	client    *http.Client
}

type jiraIssuePayload struct {
	Fields jiraFields `json:"fields"`
}

type jiraFields struct {
	Project   jiraKey `json:"project"`
	Summary   string  `json:"summary"`
	IssueType jiraKey `json:"issuetype"`
	Description string `json:"description"`
}

type jiraKey struct {
	Key string `json:"key"`
}

// NewJiraAlerter constructs a JiraAlerter from the provided config values.
func NewJiraAlerter(baseURL, token, project, issueType string) (*JiraAlerter, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("jira: baseURL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("jira: token is required")
	}
	if project == "" {
		return nil, fmt.Errorf("jira: project key is required")
	}
	if issueType == "" {
		issueType = "Task"
	}
	return &JiraAlerter{
		baseURL:   baseURL,
		token:     token,
		project:   project,
		issueType: issueType,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send creates a Jira issue describing the expiring secret.
func (j *JiraAlerter) Send(path string, expiry time.Time) error {
	payload := jiraIssuePayload{
		Fields: jiraFields{
			Project:   jiraKey{Key: j.project},
			Summary:   fmt.Sprintf("Vault secret expiring soon: %s", path),
			IssueType: jiraKey{Key: j.issueType},
			Description: fmt.Sprintf(
				"The Vault secret at path '%s' is expiring at %s. Please rotate or renew it.",
				path, expiry.UTC().Format(time.RFC3339),
			),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jira: failed to marshal payload: %w", err)
	}

	url := j.baseURL + "/rest/api/2/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jira: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+j.token)

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
