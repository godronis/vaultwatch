package alert

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// mockSNSPublisher records the last PublishInput it received.
type mockSNSPublisher struct {
	lastInput *sns.PublishInput
	errToReturn error
}

func (m *mockSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.lastInput = params
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return &sns.PublishOutput{}, nil
}

func TestNewSNSAlerter_EmptyARN(t *testing.T) {
	_, err := newSNSAlerterWithClient("", &mockSNSPublisher{})
	if err == nil {
		t.Fatal("expected error for empty topic ARN")
	}
}

func TestNewSNSAlerter_ValidARN(t *testing.T) {
	a, err := newSNSAlerterWithClient("arn:aws:sns:us-east-1:123456789012:vault-alerts", &mockSNSPublisher{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestSNSAlerter_Send_Success(t *testing.T) {
	mock := &mockSNSPublisher{}
	a, _ := newSNSAlerterWithClient("arn:aws:sns:us-east-1:123456789012:vault-alerts", mock)

	if err := a.Send("secret/myapp/db", "24h"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.lastInput == nil {
		t.Fatal("expected Publish to be called")
	}
	if !strings.Contains(*mock.lastInput.Message, "secret/myapp/db") {
		t.Errorf("expected message to contain path, got: %s", *mock.lastInput.Message)
	}
	if !strings.Contains(*mock.lastInput.Message, "24h") {
		t.Errorf("expected message to contain expiry, got: %s", *mock.lastInput.Message)
	}
	if *mock.lastInput.TopicArn != "arn:aws:sns:us-east-1:123456789012:vault-alerts" {
		t.Errorf("unexpected topic ARN: %s", *mock.lastInput.TopicArn)
	}
}

func TestSNSAlerter_Send_PublishError(t *testing.T) {
	mock := &mockSNSPublisher{errToReturn: errors.New("network timeout")}
	a, _ := newSNSAlerterWithClient("arn:aws:sns:us-east-1:123456789012:vault-alerts", mock)

	err := a.Send("secret/myapp/db", "1h")
	if err == nil {
		t.Fatal("expected error from failed publish")
	}
	if !strings.Contains(err.Error(), "publish failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
