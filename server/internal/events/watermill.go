package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"

	appConfig "github.com/aboloredev/armory/internal/config"
	"github.com/aboloredev/armory/internal/providers"
)

type PublisherEvent struct {
	publisher message.Publisher
	queueName string
}

type SubscriberEvent struct {
	subscriber message.Subscriber
	queueName string
}

func (pe *PublisherEvent) Publish(eventType string, payload interface{}, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	// ADD METADATA
	msg.Metadata.Set("event_type", eventType)
	for i, j := range metadata {
		msg.Metadata.Set(i, j)
	}

	return pe.publisher.Publish(pe.queueName, msg)
}

func (pe *PublisherEvent) Close() error {
	return pe.publisher.Close()
}

func NewEventPublisher(ctx context.Context, cfg *appConfig.AWSConfig) (*PublisherEvent, error) {
	logger := watermill.NewStdLogger(false, false)

	awsConfig, err := providers.CreateAWSConfig(ctx, &appConfig.AWSConfig{
		Region: cfg.Region, S3Endpoint: cfg.S3Endpoint, Key: cfg.Key, KeyId: cfg.KeyId})
	if err != nil {
		return nil, fmt.Errorf("Failed to create aws config %v", err)
	}

	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher %v", err)
	}

	return &PublisherEvent{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}
