package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewJiraAlerter_EmptyURL(t *testing.T) {
	_, err := NewJiraAlerter("", "token", "OPS", "Task")
	if err == nil {
		t.Fatal("expected error for empty baseURL")
	}
}

func TestNewJiraAlerter_EmptyToken(t *testing.T) {
	_, err := NewJiraAlerter("http://jira.example.com", "", "OPS", "Task")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewJiraAlerter_EmptyProject(t *testing.T) {
	_, err := NewJiraAlerter("http://jira.example.com", "token", "", "Task")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestNewJiraAlerter_DefaultIssueType(t *testing.T) {
	a, err := NewJiraAlerter("http://jira.example.com", "token", "OPS", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.issueType != "Task" {
		t.Errorf("expected default issueType 'Task', got %q", a.issueType)
	}
}

func TestNewJiraAlerter_ValidConfig(t *testing.T) {
	a, err := NewJiraAlerter("http://jira.example.com", "token", "OPS", "Bug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.project != "OPS" {
		t.Errorf("expected project 'OPS', got %q", a.project)
	}
}

func TestJiraAlerter_Send_Success(t *testing.T) {
	var captured map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	a, _ := NewJiraAlerter(server.URL, "mytoken", "OPS", "Task")
	a.client = server.Client()

	err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields, ok := captured["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("expected fields in payload")
	}
	if fields["summary"] == "" {
		t.Error("expected non-empty summary")
	}
}

func TestJiraAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	a, _ := NewJiraAlerter(server.URL, "badtoken", "OPS", "Task")
	a.client = server.Client()

	err := a.Send("secret/myapp/db", time.Now().Add(24*time.Hour))
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestJiraAlerter_Send_PayloadContainsPath(t *testing.T) {
	var captured jiraIssuePayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	a, _ := NewJiraAlerter(server.URL, "token", "SEC", "")
	a.client = server.Client()
	_ = a.Send("secret/payments/key", time.Now().Add(48*time.Hour))

	if captured.Fields.Project.Key != "SEC" {
		t.Errorf("expected project key 'SEC', got %q", captured.Fields.Project.Key)
	}
}
