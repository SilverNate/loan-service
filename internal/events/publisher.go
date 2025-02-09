package events

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (kp *KafkaPublisher) Publish(event LoanEvent) error {
	event.Timestamp = time.Now()
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(string(event.Type)),
		Value: value,
	}

	err = kp.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("error writing to kafka: %v", err)
		return err
	}

	return nil
}

//type MockPublisher struct {
//	PublishedEvents []LoanEvent
//	Err             error
//}
//
//func (mp *MockPublisher) Publish(event LoanEvent) error {
//	if mp.Err != nil {
//		return mp.Err
//	}
//	event.Timestamp = time.Unix(0, 0)
//	mp.PublishedEvents = append(mp.PublishedEvents, event)
//	return nil
//}
