package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aboloredev/armory/internal/config"
	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/logger"
	"github.com/aboloredev/armory/internal/notifications"
	"github.com/aboloredev/armory/internal/providers"
)

func init() {
	log.Println("Starting notification micro-service")
}

func main() {
	log := logger.New()
	ctx := context.Background()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load All config") 
	}

	// Initialise email notifier
	emailConfig := &notifications.SMTPConfig{
		Host: cfg.SMTP.Host,
		Port: cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From: cfg.SMTP.From,
	}

	emailNotifier := notifications.NewEmailNotifier(emailConfig)

	// The email services needs sqs to consume messages from the queue
	awsConfig, err := providers.CreateAWSConfig(ctx, &config.AWSConfig{
		Region: cfg.AWSServices.Region,
		S3Endpoint: cfg.AWSServices.S3Endpoint,
		Key: cfg.AWSServices.Key,
		KeyId: cfg.AWSServices.KeyId,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not load AWSconfig") 
	}

	// Start watermill
	logger := watermill.NewStdLogger(false, false)
	
	subscriberConfig := sqs.SubscriberConfig{
		AWSConfig: awsConfig,
	}

	subscriber, err := sqs.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		log.Error().Err(err).Msg("Subscriber failed to create") 
	}

	message, err := subscriber.Subscribe(ctx, cfg.AWSServices.EventQueueName)
	if err != nil {
		subscriber.Close()
		log.Fatal().Err(err).Msg("Failed to subscribe to queue") 
	}

	// Initiate graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("notification micro-service has started and waiting for signal")

	for {
		select {
		case msg := <-message: 
			ProcessMessage(msg, emailNotifier)
			if err != nil {
				log.Printf("Error processing message: %s", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigChan:
			log.Println("notification micro-service shuttin down")
			subscriber.Close()
			return
		}
	}
}

func ProcessMessage(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case notifications.OrderCreatedSuccessfully:
		return HandleOrderCreated(msg, emailNotifier)
	default: 
		log.Printf("Unknown event type %s", eventType)
		return nil
	} 
}

func HandleOrderCreated(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	var order *dto.OrderResponse
	fmt.Println(order)

	err := json.Unmarshal(msg.Payload, &order)
	if err != nil {
		return err
	}

	userName := order.UserFirstName + " " + order.UserLastName
	if userName == "" {
		userName = "user"
	}

	userEmail := order.UserEmail
	if userEmail == "" {
		userEmail = "userEmail"
	}

	userOrder := order

	log.Printf("Sending order notification to %s", userName)

	emailNotifier.SendOrderNotificationMail(userEmail, userName, userOrder)

	log.Println("Notification sent")

	return nil
}
