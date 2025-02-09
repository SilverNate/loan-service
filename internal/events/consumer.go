package events

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"loan-service/internal/logger"
	"time"
)

type EventConsumer struct {
	reader *kafka.Reader
	logger logger.Logger
}

func NewKafkaConsumer(brokers []string, topic, groupID string, logger logger.Logger) *EventConsumer {
	return &EventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:          brokers,
			GroupID:          groupID,
			Topic:            topic,
			MinBytes:         10e3, // 10KB
			MaxBytes:         10e6, // 10MB
			RebalanceTimeout: 10 * time.Second,
			JoinGroupBackoff: 10 * time.Second,
		}),
		logger: logger,
	}
}

func (ec *EventConsumer) Consume(ctx context.Context) error {
	for {
		message, err := ec.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var event LoanEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			ec.logger.Error("Error unmarshalling event: ", err)
			continue
		}
		ec.ProcessEvent(event)

	}
}

func (ec *EventConsumer) ProcessEvent(event LoanEvent) {
	switch event.Type {
	case LoanCreated:
		ec.logger.Infof("Processed LoanCreated event for loan id: %d", event.LoanID)
		// TODO: we can add trigger notification to firebase
	case LoanApproved:
		ec.logger.Infof("Processed LoanApproved event for loan id: %d", event.LoanID)
		// TODO: we can add trigger notification to firebase
	case LoanInvested:
		ec.logger.Infof("Processed LoanInvested event for loan id: %d", event.LoanID)
		// TODO: we can add trigger notification to firebase
		// give user notification ROI
	case LoanDisbursed:
		ec.logger.Infof("Processed LoanDisbursed event for loan id: %d", event.LoanID)
		// TODO: we can add trigger notification to firebase
	default:
		ec.logger.Infof("Received unknown event type: %s for loan id: %d", event.Type, event.LoanID)
	}
}
