package alert

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSPublisher abstracts the SNS Publish call for testability.
type SNSPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSAlerter sends Vault expiration alerts to an AWS SNS topic.
type SNSAlerter struct {
	topicARN string
	client   SNSPublisher
}

// NewSNSAlerter creates a new SNSAlerter. It loads AWS credentials from the
// default credential chain (env vars, ~/.aws/credentials, IAM role, etc.).
func NewSNSAlerter(topicARN string) (*SNSAlerter, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: failed to load AWS config: %w", err)
	}

	return &SNSAlerter{
		topicARN: topicARN,
		client:   sns.NewFromConfig(cfg),
	}, nil
}

// newSNSAlerterWithClient creates an SNSAlerter with a pre-built publisher
// (used in tests).
func newSNSAlerterWithClient(topicARN string, pub SNSPublisher) (*SNSAlerter, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	return &SNSAlerter{topicARN: topicARN, client: pub}, nil
}

// Send publishes a secret-expiration alert message to the configured SNS topic.
func (a *SNSAlerter) Send(path string, expiresIn string) error {
	msg := fmt.Sprintf("VaultWatch Alert: secret '%s' expires in %s", path, expiresIn)

	_, err := a.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(a.topicARN),
		Message:  aws.String(msg),
		Subject:  aws.String("Vault Secret Expiring Soon"),
	})
	if err != nil {
		return fmt.Errorf("sns: publish failed: %w", err)
	}
	return nil
}
