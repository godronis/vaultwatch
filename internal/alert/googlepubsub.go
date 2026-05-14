package alert

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// pubsubClient abstracts the Google Pub/Sub topic for testing.
type pubsubClient interface {
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
}

// GooglePubSubAlerter sends alerts to a Google Cloud Pub/Sub topic.
type GooglePubSubAlerter struct {
	topic  pubsubClient
	topicID string
}

// NewGooglePubSubAlerter creates a new GooglePubSubAlerter.
// projectID and topicID are required.
func NewGooglePubSubAlerter(projectID, topicID string) (*GooglePubSubAlerter, error) {
	if projectID == "" {
		return nil, fmt.Errorf("googlepubsub: projectID is required")
	}
	if topicID == "" {
		return nil, fmt.Errorf("googlepubsub: topicID is required")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("googlepubsub: failed to create client: %w", err)
	}

	return &GooglePubSubAlerter{
		topic:   client.Topic(topicID),
		topicID: topicID,
	}, nil
}

// newGooglePubSubAlerterWithClient creates an alerter with an injected client (for testing).
func newGooglePubSubAlerterWithClient(topic pubsubClient, topicID string) *GooglePubSubAlerter {
	return &GooglePubSubAlerter{topic: topic, topicID: topicID}
}

// Send publishes an alert message to the configured Pub/Sub topic.
func (a *GooglePubSubAlerter) Send(path string, expiresIn string) error {
	payload := map[string]string{
		"path":       path,
		"expires_in": expiresIn,
		"alert":      fmt.Sprintf("Secret at %s is expiring in %s", path, expiresIn),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlepubsub: failed to marshal payload: %w", err)
	}

	ctx := context.Background()
	result := a.topic.Publish(ctx, &pubsub.Message{Data: data})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("googlepubsub: failed to publish message: %w", err)
	}

	return nil
}
