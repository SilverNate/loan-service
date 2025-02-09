package main

import (
	"context"
	"log"
	"os"
	"time"

	"loan-service/internal/events"
	"loan-service/internal/logger"
)

func main() {
	logrusLogger := logger.NewLogrusLogger()

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable is not set")
	}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "loan_events"
	}
	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "loan_consumer_group"
	}

	consumer := events.NewKafkaConsumer([]string{kafkaBrokers}, kafkaTopic, groupID, logrusLogger)

	ctx := context.Background()
	for {
		if err := consumer.Consume(ctx); err != nil {
			logrusLogger.Errorf("Error consuming events: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}
